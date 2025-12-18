/**
 * Auto-assign PR to commit authors
 *
 * This script automatically assigns a pull request to all users who have
 * authored commits in the PR, excluding bot users.
 *
 * @param {Object} github - GitHub API client from actions/github-script
 * @param {Object} context - Workflow context from actions/github-script
 */
module.exports = async ({ github, context }) => {
  try {
    // Fetch PR commits and current assignees in parallel
    const [{ data: prCommits }, { data: { assignees } }] = await Promise.all([
      github.rest.pulls.listCommits({
        owner: context.repo.owner,
        repo: context.repo.repo,
        pull_number: context.issue.number
      }),
      github.rest.pulls.get({
        owner: context.repo.owner,
        repo: context.repo.repo,
        pull_number: context.issue.number
      })
    ]);

    // Extract unique commit authors (excluding bots)
    const commitAuthors = extractCommitAuthors(prCommits);

    // Calculate new assignees (authors not already assigned)
    const existingAssignees = new Set(assignees.map(a => a.login));
    const newAssignees = Array.from(commitAuthors).filter(
      author => !existingAssignees.has(author)
    );

    // Add new assignees if any
    if (newAssignees.length > 0) {
      await addAssigneesToPR(github, context, newAssignees);
    } else {
      console.log('ℹ️ No new assignees to add');
    }
  } catch (error) {
    console.error(`❌ Failed to process PR assignees: ${error.message}`);
    throw error;
  }
};

/**
 * Extract GitHub usernames from commits, excluding bots
 *
 * @param {Array} commits - Array of commit objects from GitHub API
 * @returns {Set<string>} Set of unique GitHub usernames
 */
function extractCommitAuthors(commits) {
  const authors = new Set();

  for (const commit of commits) {
    if (commit.author?.login && commit.author.type !== 'Bot') {
      authors.add(commit.author.login);
    } else if (commit.author?.type === 'Bot') {
      console.log(`⚠️ Skipping bot commit ${commit.sha.substring(0, 7)}`);
    } else {
      console.log(`⚠️ Skipping commit ${commit.sha.substring(0, 7)} (No linked GitHub user)`);
    }
  }

  return authors;
}

/**
 * Add assignees to a pull request
 *
 * @param {Object} github - GitHub API client
 * @param {Object} context - Workflow context
 * @param {Array<string>} assignees - Array of GitHub usernames to assign
 */
async function addAssigneesToPR(github, context, assignees) {
  try {
    console.log(`✅ Adding assignees: ${assignees.join(', ')}`);

    await github.rest.issues.addAssignees({
      owner: context.repo.owner,
      repo: context.repo.repo,
      issue_number: context.issue.number,
      assignees: assignees
    });
  } catch (error) {
    console.error(`❌ Failed to add assignees: ${error.message}`);
    throw error;
  }
}
