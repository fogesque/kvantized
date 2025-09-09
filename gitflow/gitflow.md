# Git Workflow Guidelines

This document outlines the standard Git workflow for our development team. Adhering to these guidelines ensures a consistent, clean, and collaborative development process.

---

## 1. Branching Strategy

We will follow a branching strategy inspired by Git Flow, with specific naming conventions and restrictions.

**Key Principles:**

* **No Direct Pushes to `main` or `develop`:** The `main` and `develop` branches are protected and should only be updated via merged Pull Requests.
* **Feature Branches:** For new features or significant changes.
* **Hotfix Branches:** For urgent bug fixes in production.
* **Release Branches:** For preparing a new production release.

**Branch Naming Conventions:**

All branches should start with a specific prefix followed by a descriptive, lowercase, hyphen-separated name.

* **Feature Branches:**
    * **Prefix:** `feature/`
    * **Format:** `feature/<issue-id>-<short-description>`
    * **Examples:** `feature/kanon-69-network-module`, `feature/kanon-666-api-endpoint-for-control`
    * **Purpose:** Developed from `develop`. Merged back into `develop`.

* **Hotfix Branches:**
    * **Prefix:** `hotfix/`
    * **Format:** `hotfix/<issue-id>-<short-description>`
    * **Examples:** `hotfix/kanon-69-fix-conn-error`, `hotfix/kanon-666-patch-security-vulnerability`
    * **Purpose:** Developed from `develop`. Merged back into `develop`.

* **Temporary/Experiment Branches:**
    * **Prefix:** `tmp/` or `exp/`
    * **Format:** `tmp/<your-name>-<description>`
    * **Examples:** `tmp/john-test-branch`, `exp/jane-api-experiment`
    * **Purpose:** For quick tests or experiments. These branches should be short-lived and deleted once their purpose is served.

---

## 2. General Workflow Steps

1.  **Pull Latest Changes:** Always start by pulling the latest changes from the `develop` branch to ensure your local `develop` is up-to-date.
    ```bash
    git checkout develop
    git pull origin develop
    ```

2.  **Create a New Branch:** Based on the type of work you're doing, create a new branch from `develop` (for features) or `main` (for hotfixes).
    ```bash
    # For a new feature
    git checkout -b feature/<issue-id>-<your-feature-description> develop

    # For a hotfix
    git checkout -b hotfix/<issue-id>-<your-hotfix-description> main
    ```

3.  **Develop and Commit Regularly:** Make small, logical commits with clear and concise commit messages.

    * **Commit Message Format:** `[JIRA tag]: Verb Subject`
        * **Verb:** `Add` (new feature), `Fix` (bug fix), `Document` (documentation), `Format` (formatting, no code changes), `Refactor` (refactoring production code), `Test` (adding tests)
        * **JIRA tag (Optional):** JIRA task tag
        * **Subject:** Concise description of the change.

    * **Examples:**
        * `KANON-69: Add driver loading method`
        * `KANON-666: Fix driver loading error`

4.  **Run Formatter and Linter:** Before committing, always run the configured code formatter and linter. This ensures code consistency and catches potential issues early.

    * **Clang-Format:**
        * Configured in VS Code codestyle.profile

    * **Clang-Tidy**
        * Configured in VS Code codestyle.profile

5.  **Push Your Branch:** Regularly push your branch to GitLab, especially when you have significant progress or before leaving for the day.
    ```bash
    git push origin <your-branch-name>
    ```

---

## 3. Pull Request (PR) Process (GitLab)

Once your feature or hotfix is complete and thoroughly tested on your local machine, it's time to open a Pull Request.

1.  **Rebase:** Before creating a PR, consider rebasing your feature branch onto the latest `develop`. This keeps your commit history linear and clean.
    ```bash
    # From your feature branch
    git checkout develop
    git pull
    git checkout <your_branch>
    git rebase develop
    # Resolve any conflicts if they occur
    git push --force-with-delete origin <your-branch-name> # USE WITH CAUTION! Only if you understand rebase and its implications.
    ```
    **Note on `git push --force-with-delete`:** Only use this if you are the *only* person working on the branch and you understand the implications of rewriting history. For shared branches, prefer merging `develop` into your feature branch.

2.  **Navigate to GitLab:** Go to your project on GitLab.

3.  **Create New Pull Request:**
    * GitLab will often show a prompt to create a new merge request when you push a new branch.
    * Alternatively, go to "Merge Requests" on the left sidebar and click "New merge request".

4.  **Configure the Pull Request:**
    * **Source Branch:** Select your branch (`feature/<issue-id>-<your-feature-description>` or `hotfix/<issue-id>-<your-hotfix-description>`).
    * **Target Branch:**
        * For `feature/` branches: `develop`
        * For `hotfix/` branches: `main`

5.  **Title and Description:**
    * **Title:** A concise summary of the PR's purpose (often similar to your main commit message).
    * **Description:**
        * Clearly explain what the PR does, why it's needed, and how it addresses the associated issue.
        * Reference the relevant issue tracker ID (e.g., `#KANON-69`).
        * Include any necessary context, screenshots, or videos.
        * Mention any specific areas to review.

6.  **Assign Reviewers:** Assign at least one reviewer from the team.

7.  **Add Labels:** Add relevant labels (e.g., `feature`, `bug`, `frontend`, `backend`).

8.  **Submit Pull Request:** Click "Create merge request."

---

## 4. Code Review Process

1.  **Respond to Feedback:** Be responsive to comments and suggestions from reviewers. Engage in constructive discussions.
2.  **Make Changes:** Make necessary changes on your branch, commit them, and push them. The PR will automatically update.
3.  **Address All Comments:** Ensure all comments are addressed or clarified before the PR is merged.
4.  **Approve and Merge:** Once the PR has received the required approvals and all checks (CI/CD) pass, the reviewer will merge the PR.

---

## 5. Post-Merge

1.  **Delete Branch:** After your PR is merged, delete your local branch and the remote branch (GitLab usually offers this option automatically when merging).
    ```bash
    git checkout develop # Or main, depending on where you merged
    git pull origin develop # Or main
    git branch -d <your-branch-name> # Delete local branch
    git push origin --delete <your-branch-name> # Delete remote branch (if not deleted by GitLab)
    ```

2.  **Pull Latest `develop`:** Always pull the latest `develop` (or `main`) after a PR has been merged to keep your local environment synchronized.

---

By following these guidelines, we can maintain a clean, organized, and efficient Git repository.