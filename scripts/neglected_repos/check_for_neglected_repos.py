import os
import requests
from datetime import datetime, timedelta


def get_github_headers():
  """Retrieve GitHub headers with authorization."""
  github_token = os.getenv("GH_TOKEN")
  return {"Authorization": f"token {github_token}"}


def get_main_branch_last_commit(owner, repo, headers):
  """Fetch the last commit date of the main branch."""
  base_url = "https://api.github.com"
  commits_url = f"{base_url}/repos/{owner}/{repo}/commits/main"

  commits_response = requests.get(commits_url, headers=headers)
  commits_response.raise_for_status()
  commit_data = commits_response.json()

  last_commit_date_str = commit_data["commit"]["committer"]["date"]
  last_commit_date = datetime.fromisoformat(last_commit_date_str.replace("Z", "+00:00"))

  print(f"Last commit to 'main':")
  print(f"  SHA: {commit_data['sha']}")
  print(f"  Date: {last_commit_date}")

  return last_commit_date


def calculate_time_periods(reference_date):
  """Calculate key time periods based on the reference date."""
  now = datetime.now(reference_date.tzinfo)
  return {
    "two_weeks_ago": now - timedelta(weeks=2),
    "three_months_ago": now - timedelta(weeks=12),
    "one_year_ago": now - timedelta(weeks=52),
  }


def get_open_pull_requests(owner, repo, headers):
  """Fetch open pull requests for the repository."""
  base_url = "https://api.github.com"
  pulls_url = f"{base_url}/repos/{owner}/{repo}/pulls?state=open"

  pulls_response = requests.get(pulls_url, headers=headers)
  pulls_response.raise_for_status()
  pulls_data = pulls_response.json()

  print(f"\nOpen Pull Requests: {len(pulls_data)}")
  return pulls_data


def analyse_branches(owner, repo, pulls_data, headers, time_periods):
  """Analyze branches for commit activity and maintenance status."""
  base_url = "https://api.github.com"
  active_maintenance_pr = False
  commit_three_months_ago = False
  commit_two_weeks_old = False

  for pr in pulls_data:
    branch_name = pr["head"]["ref"]
    print(f"  Branch: {branch_name}")

    branch_url = f"{base_url}/repos/{owner}/{repo}/branches/{branch_name}"
    branch_response = requests.get(branch_url, headers=headers)
    branch_response.raise_for_status()
    branch_data = branch_response.json()

    branch_last_commit_date_str = branch_data["commit"]["commit"]["committer"]["date"]
    branch_last_commit_date = datetime.fromisoformat(branch_last_commit_date_str.replace("Z", "+00:00"))
    print(f"    Last Commit Date: {branch_last_commit_date}")

    if branch_last_commit_date < time_periods["three_months_ago"]:
      print("    The last commit was more than three months ago.")
      commit_three_months_ago = True
      commit_two_weeks_old = True
    elif branch_last_commit_date < time_periods["two_weeks_ago"]:
      print("    The last commit was more than two weeks ago.")
      commit_two_weeks_old = True
    else:
      print("    The last commit was within the past two weeks.")

    if branch_name.startswith("renovate") or branch_name.startswith("dependabot"):
      print("    Renovate PR.")
      if commit_two_weeks_old:
        active_maintenance_pr = True

  return active_maintenance_pr, commit_three_months_ago, commit_two_weeks_old


def generate_message(repo, owner, main_commit_one_year_ago, active_maintenance_pr, commit_three_months_ago):
  message = ""

  if active_maintenance_pr or commit_three_months_ago:
    link = f"https://github.com/{owner}/{repo}/pulls"
  else:
    link = f"https://github.com/{owner}/{repo}"

  if main_commit_one_year_ago or active_maintenance_pr or commit_three_months_ago:
    message += f"Repository *{repo}* needs action against it because:\n"
    message += f"{link}\n"
    if main_commit_one_year_ago:
      message += "   - No commits in over a year (consider archiving or updating)\n"
    if active_maintenance_pr:
      message += "   - Dependency maintenance PR (renovate/dependabot) open for more than 2 weeks (merge or close the PR)\n"
    if commit_three_months_ago:
      message += "   - Stale PR has remained open for more than 3 months (merge or close the PR)\n"

  return message


def fetch_github_repo_data(owner, repo):
  """Determines the state of the repository based on commit and PR activity."""
  try:
    headers = get_github_headers()
    last_commit_date = get_main_branch_last_commit(owner, repo, headers)
    time_periods = calculate_time_periods(last_commit_date)

    main_commit_one_year_ago = last_commit_date < time_periods["one_year_ago"]
    pulls_data = get_open_pull_requests(owner, repo, headers)
    active_maintenance_pr, commit_three_months_ago, commit_two_weeks_old = analyse_branches(
      owner, repo, pulls_data, headers, time_periods
    )

    # Generate final output message
    if main_commit_one_year_ago or active_maintenance_pr or commit_three_months_ago:
      output_message = generate_message(
        repo,
        owner,
        main_commit_one_year_ago,
        active_maintenance_pr,
        commit_three_months_ago,
      )
      return output_message

    return ""
  except requests.exceptions.RequestException as e:
    print(f"An error occurred: {e}")
    return ""


def fetch_team_repos(org, team):
  """Fetch all repositories belonging to a specific team within an organization"""
  base_url = f"https://api.github.com/orgs/{org}/teams/{team}/repos"
  headers = get_github_headers()
  matching_repos = []

  try:
    page = 1
    while True:
      response = requests.get(
        base_url,
        headers=headers,
        params={"page": page, "per_page": 50}
      )
      response.raise_for_status()
      repos = response.json()

      if not repos:
        break

      for repo in repos:
        # Ensure the repository belongs to the specified org
        if repo['owner']['login'].lower() == org.lower() and not repo['archived']:
          matching_repos.append(repo['name'])

      page += 1

    print(f"\nRepositories in team '{team}':")
    for repo_name in matching_repos:
      print(f"  - {repo_name}")

    return matching_repos
  except requests.exceptions.RequestException as e:
    print(f"An error occurred while fetching team repositories: {e}")
    return []


def send_slack_message(messages):
  """
  Sends a formatted list of messages to Slack using a webhook.
  """
  webhook_url = os.getenv('SLACK_WEBHOOK')
  if not webhook_url:
    print("SLACK_WEBHOOK environment variable is not set.")
    exit(1)

  title = ":workshop: *--- Repositories in need of action ---* :workshop:\n\n"
  formatted_message = "\n\n".join(f"{message}" for message in messages)
  final_message = title + formatted_message

  if len(final_message) > 40000:
    print("Message exceeds Slack's 40,000 character limit. Please reduce message size.")
    exit(1)

  payload = {
    "text": final_message
  }

  # Send the message to Slack
  response = requests.post(webhook_url, json=payload)

  if response.status_code != 200:
    print(f"Failed to send message to Slack: {response.status_code}, {response.text}")
    exit(1)
  else:
    print("Slack message sent successfully")

  return response


if __name__ == "__main__":
  owner = "ministryofjustice"
  team = "opg"

  repos_to_check = fetch_team_repos(owner, team)

  output_messages = []
  for repo in repos_to_check:
    output_message = fetch_github_repo_data(owner, repo)
    if len(output_message) > 0:
      output_messages.append(output_message)

  for msg in output_messages:
    print(msg)

  if len(output_messages) > 0:
    print("-- Preparing to Send Message --")
    send_slack_message(output_messages)
