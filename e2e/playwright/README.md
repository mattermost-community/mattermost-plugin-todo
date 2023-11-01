In order to get your environment set up to run [Playwright](https://playwright.dev) tests, you can run `./setup-environment`, or run equivalent commands for your current setup.

What this script does:

- Navigate to the folder above `mattermost-plugin-todo`
- Clone `mattermost` (if it is already cloned there, please have a clean git index to avoid issues with conflicts)
- `cd mattermost`
- Install webapp dependencies - `cd webapp && npm i`
- Install Playwright test dependencies - `cd ../e2e-tests/playwright && npm i`
- Install Playwright - `npx install playwright`
- Install Todo plugin e2e dependencies - `cd ../../../mattermost-plugin-todo/e2e/playwright && npm i`
- Build and deploy plugin with e2e support - `make deploy-e2e`

---

Then to run the tests:

Start Mattermost server:

- `cd <path>/mattermost/server`
- `make test-data`
- `make run-server`

Run test:

- `cd <path>/mattermost-plugin-todo/e2e/playwright`
- `npm test`

To see the test report:

- `cd <path>/mattermost-plugin-todo/e2e/playwright`
- `npm run show-report`
- Navigate to http://localhost:9323

To see test screenshots:

- `cd <path>/mattermost-plugin-todo/e2e/playwright/screenshots`
