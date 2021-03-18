package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"runtime/pprof"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/gui"
	"github.com/coyim/coyim/gui/settings"
	"github.com/coyim/gotk3adapter/gdk_mock"
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtk_mock"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pango_mock"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	. "gopkg.in/check.v1"
)

type MainSuite struct{}

var _ = Suite(&MainSuite{})

func (s *MainSuite) Test_mainInit_setsVersionFromCommit(c *C) {
	origBuildTag := BuildTag
	origBuildCommit := BuildCommit
	origCoyimVersion := coyimVersion

	defer func() {
		BuildTag = origBuildTag
		BuildCommit = origBuildCommit
		coyimVersion = origCoyimVersion
	}()

	BuildTag = "(no tag)"
	BuildCommit = "hello abc"
	mainInit()
	c.Assert(coyimVersion, Equals, "hello abc")

	BuildTag = ""
	BuildCommit = "hello def"
	mainInit()
	c.Assert(coyimVersion, Equals, "hello def")

	BuildTag = "v42"
	BuildCommit = "hello def"
	mainInit()
	c.Assert(coyimVersion, Equals, "v42")
}

type mockLooper struct {
	called bool
}

func (m *mockLooper) Loop() {
	m.called = true
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func captureStderr(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func captureStdoutAndStderr(f func()) (stdout, stderr string) {
	stdout = captureStdout(func() {
		stderr = captureStderr(f)
	})

	return
}

func (s *MainSuite) Test_main_parsesFlagsAndRunsClient(c *C) {
	ll := log.StandardLogger()
	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.PanicLevel)

	orgCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = orgCommandLine
	}()

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine = fs

	orgCreateGTK := createGTK
	defer func() {
		createGTK = orgCreateGTK
	}()

	called1 := false
	ml := &mockLooper{}
	createGTK = func(g gui.Graphics) looper {
		called1 = true
		return ml
	}

	stdout, _ := captureStdoutAndStderr(main)

	c.Assert(stdout, Equals, "\n")
	c.Assert(fs.Parsed(), Equals, true)
	c.Assert(called1, Equals, true)
	c.Assert(ml.called, Equals, true)
}

type gtkMockWithSettings struct {
	*gtk_mock.Mock
}

func (*gtkMockWithSettings) ApplicationNew(appID string, flags glibi.ApplicationFlags) (gtki.Application, error) {
	return &gtk_mock.MockApplication{}, nil
}

func (*gtkMockWithSettings) SettingsGetDefault() (gtki.Settings, error) {
	return &gtk_mock.MockSettings{}, nil
}

type mockGlib struct {
	glib_mock.Mock

	schemaSourceToReturn  glibi.SettingsSchemaSource
	settingsNewFullReturn glibi.Settings
}

type mockSimpleAction struct {
	*glib_mock.MockObject
}

func (*mockSimpleAction) GetName() string                     { return "" }
func (*mockSimpleAction) GetEnabled() bool                    { return false }
func (*mockSimpleAction) GetState() glibi.Variant             { return nil }
func (*mockSimpleAction) GetStateHint() glibi.Variant         { return nil }
func (*mockSimpleAction) GetParameterType() glibi.VariantType { return nil }
func (*mockSimpleAction) GetStateType() glibi.VariantType     { return nil }
func (*mockSimpleAction) ChangeState(value glibi.Variant)     {}
func (*mockSimpleAction) Activate(parameter glibi.Variant)    {}
func (*mockSimpleAction) SetEnabled(enabled bool)             {}
func (*mockSimpleAction) SetState(value glibi.Variant)        {}

func (m *mockGlib) SimpleActionNew(name string, parameterType glibi.VariantType) glibi.SimpleAction {
	return &mockSimpleAction{}
}

func (m *mockGlib) SettingsSchemaSourceNewFromDirectory(v1 string, v2 glibi.SettingsSchemaSource, v3 bool) glibi.SettingsSchemaSource {
	return m.schemaSourceToReturn
}

func (m *mockGlib) SettingsNewFull(v1 glibi.SettingsSchema, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	return m.settingsNewFullReturn
}

type mockSettingsSchema struct {
	glib_mock.MockSettingsSchema
}

type mockSettingsSchemaSource struct {
	glib_mock.MockSettingsSchemaSource

	lookupReturn glibi.SettingsSchema
}

func (m *mockSettingsSchemaSource) Lookup(v1 string, v2 bool) glibi.SettingsSchema {
	return m.lookupReturn
}

func (s *MainSuite) Test_createGTK_works(c *C) {
	oldEnv := os.Getenv("XDG_DATA_HOME")
	defer func() {
		os.Setenv("XDG_DATA_HOME", oldEnv)
	}()
	os.Setenv("XDG_DATA_HOME", "somewhere-hopefully-not-existing")

	orgHooks := hooks
	defer func() {
		hooks = orgHooks
	}()

	hooks = func() gui.OSHooks {
		return &gui.NoHooks{}
	}

	ss := &mockSettingsSchemaSource{}
	sch := &mockSettingsSchema{}
	ss.lookupReturn = sch

	mg := &mockGlib{
		schemaSourceToReturn: ss,
	}

	settings.InitSettings(mg)

	res := createGTK(gui.CreateGraphics(
		&gtkMockWithSettings{},
		mg,
		&gdk_mock.Mock{},
		&pango_mock.Mock{},
	))
	c.Assert(res, Not(IsNil))
}

func (s *MainSuite) Test_main_printsVersionAndQuits(c *C) {
	orgCommandLine := flag.CommandLine
	defer func() {
		flag.CommandLine = orgCommandLine
	}()

	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine = fs

	orgVersionFlag := *config.VersionFlag
	defer func() {
		*config.VersionFlag = orgVersionFlag
	}()

	*config.VersionFlag = true

	stdout, _ := captureStdoutAndStderr(main)

	c.Assert(stdout, Matches, "CoyIM version .*\n")
}

func (s *MainSuite) Test_initLog_setsDebugFlag(c *C) {
	ll := log.StandardLogger()
	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.PanicLevel)

	orgDebugFlag := *config.DebugFlag
	defer func() {
		*config.DebugFlag = orgDebugFlag
	}()

	*config.DebugFlag = true

	initLog()

	c.Assert(ll.Level, Equals, log.DebugLevel)
}

func (s *MainSuite) Test_initLog_setsTraceFlag(c *C) {
	ll := log.StandardLogger()
	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.PanicLevel)

	orgDebugFlag := *config.DebugFlag
	defer func() {
		*config.DebugFlag = orgDebugFlag
	}()

	*config.DebugFlag = true

	orgTraceFlag := *config.TraceFlag
	defer func() {
		*config.TraceFlag = orgTraceFlag
	}()

	*config.TraceFlag = true

	initLog()

	c.Assert(ll.Level, Equals, log.TraceLevel)
}

func (s *MainSuite) Test_stopProfileIfNecessary_works(c *C) {
	orgCPUProfileFlag := *config.CPUProfile
	defer func() {
		*config.CPUProfile = orgCPUProfileFlag
	}()
	*config.CPUProfile = ""

	stopProfileIfNecessary()

	*config.CPUProfile = "somewhere-not-existing-hopefully"

	stopProfileIfNecessary()
}

func (s *MainSuite) Test_startProfileIfNecessary_failsOnBadFile(c *C) {
	defer pprof.StopCPUProfile()

	orgCPUProfileFlag := *config.CPUProfile
	defer func() {
		*config.CPUProfile = orgCPUProfileFlag
	}()
	*config.CPUProfile = "somewhere/that/hopefully/doesnt/exists"

	ll := log.StandardLogger()
	orgExitFunc := ll.ExitFunc
	defer func() {
		ll.ExitFunc = orgExitFunc
	}()

	called := false
	ll.ExitFunc = func(int) {
		called = true
	}

	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.DebugLevel)
	hook := test.NewGlobal()

	startProfileIfNecessary()
	c.Assert(called, Equals, true)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.FatalLevel)
	c.Assert(hook.LastEntry().Message, Matches, "could not create CPU profile")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, ".*(no such file or directory|cannot find the path specified).*")
}

func (s *MainSuite) Test_startProfileIfNecessary_failsOnStarting(c *C) {
	defer pprof.StopCPUProfile()

	orgCPUProfileFlag := *config.CPUProfile
	defer func() {
		*config.CPUProfile = orgCPUProfileFlag
	}()
	*config.CPUProfile = "somewhere/that/hopefully/doesnt/exists"

	ll := log.StandardLogger()
	orgExitFunc := ll.ExitFunc
	defer func() {
		ll.ExitFunc = orgExitFunc
	}()

	called := false
	ll.ExitFunc = func(int) {
		called = true
	}

	orgLevel := ll.Level
	defer func() {
		ll.SetLevel(orgLevel)
	}()

	ll.SetLevel(log.DebugLevel)
	hook := test.NewGlobal()

	orgOsCreate := osCreate
	defer func() {
		osCreate = orgOsCreate
	}()

	osCreate = func(string) (*os.File, error) {
		return nil, nil
	}

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	pprof.StartCPUProfile(tmpfile)

	startProfileIfNecessary()
	c.Assert(called, Equals, true)

	c.Assert(len(hook.Entries), Equals, 1)
	c.Assert(hook.LastEntry().Level, Equals, log.FatalLevel)
	c.Assert(hook.LastEntry().Message, Matches, "could not start CPU profile")
	c.Assert(hook.LastEntry().Data["error"], ErrorMatches, "cpu profiling already in use")
}

func (s *MainSuite) Test_startProfileIfNecessary_works(c *C) {
	defer pprof.StopCPUProfile()

	tmpfile, _ := ioutil.TempFile("", "")
	defer os.Remove(tmpfile.Name())

	orgCPUProfileFlag := *config.CPUProfile
	defer func() {
		*config.CPUProfile = orgCPUProfileFlag
	}()
	*config.CPUProfile = tmpfile.Name()

	orgOsCreate := osCreate
	defer func() {
		osCreate = orgOsCreate
	}()

	osCreate = func(string) (*os.File, error) {
		return tmpfile, nil
	}

	startProfileIfNecessary()
}
