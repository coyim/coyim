const core = require('@actions/core');

const colors = {
    'success': '#00FF2D',
    'cancelled': '#FFF608',
    'failure': '#FF0A00',
};

const statusEndings = {
    'success': 'was successful!',
    'cancelled': 'was cancelled!',
    'failure': 'failed!',
};

function colorFor(status) {
    return colors[status.toLowerCase()] || '#000000';
}

function statusEndingFor(status) {
    return statusEndings[status.toLowerCase()] || `has status '${status}'`;
}

function calculateStatus(status, commitID, commitMessage) {
    const runId = process.env.GITHUB_RUN_ID || '<unknown>';
    const runNumber = process.env.GITHUB_RUN_NUMBER || '<unknown>';
    const githubServer = process.env.GITHUB_SERVER_URL || '<unknown>';
    const repo = process.env.GITHUB_REPOSITORY || '<unknown>';
    const buildURL = `${githubServer}/${repo}/actions/runs/${runId}`;
    const workflow = process.env.GITHUB_WORKFLOW || '<unknown>';
    const actor = process.env.GITHUB_ACTOR || '<unknown>';
    const actorURL = `${githubServer}/${actor}`;

    const subject = `${status.toUpperCase()} Build #${runNumber} received status ${status}!`;

    let message = `<h1><span data-mx-color="${colorFor(status)}">${status.toUpperCase()}</span></h1>`;
    message += `Build <a href="${buildURL}"> ${repo} #${runNumber} ${workflow}</a> `;
    message += statusEndingFor(status);
    message += `<br>Triggered by <a href="${actorURL}">${actor}</a>`;
    if (commitMessage) {
        message += `: <i>${commitMessage}</i>`;
    }
    if (commitID) {
        const shortCommitID = commitID.substring(0, 8);
        const commitURL = `${githubServer}/${repo}/commit/${commitID}`;

        message += ` (<a href="${commitURL}">#${shortCommitID}</a>)`;
    }

    core.setOutput('subject', subject);
    core.setOutput('message', message);
}

async function run() {
    try {
        const status = core.getInput('status');
        const commitID = core.getInput('commit_id');
        const commitMessage = core.getInput('commit_message');
        calculateStatus(status, commitID, commitMessage);
    }
    catch (error) {
        core.setFailed(error.message);
    }
}

run();

