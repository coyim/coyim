const core = require('@actions/core');

const colors = {
    'success': '#00FF2D',
    'cancelled': '#FFF608',
    'failure': '#FF0A00',
};

const statusEndings = {
    'success': 'was successful!',
    'cancelled': 'failed!',
    'failure': 'was cancelled!',
};

function colorFor(status) {
    return colors[status.toLowerCase()] || '#000000';
}

function statusEndingFor(status) {
    return statusEndings[status.toLowerCase()] || `has status '${status}'`;
}

function calculateStatus(status) {
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
    message += `<br>Triggered by <a href="${actorURL}">${actor}</a> `;

    core.setOutput('subject', subject);
    core.setOutput('message', message);
}

async function run() {
    try {
        const status = core.getInput('status');
        calculateStatus(status);
    }
    catch (error) {
        core.setFailed(error.message);
    }
}

run();

