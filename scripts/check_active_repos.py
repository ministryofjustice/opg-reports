import os
import requests


def fetch_github_repo_data(owner, repo):
  # Base URL for GitHub API
  base_url = "https://api.github.com"

  # Endpoint to get the main branch's latest commit
  commits_url = f"{base_url}/repos/{owner}/{repo}/commits/main"

  # Endpoint to get the open pull requests
  pulls_url = f"{base_url}/repos/{owner}/{repo}/pulls?state=open"

  github_token = os.getenv('GH_TOKEN')
  headers = {"Authorization": f"token {github_token}"}

  try:
    # Fetch the latest commit details
    commits_response = requests.get(commits_url, headers=headers)
    commits_response.raise_for_status()  # Raise an error if the request fails
    commit_data = commits_response.json()

    # Extract the last commit information
    last_commit_sha = commit_data['sha']
    last_commit_date = commit_data['commit']['committer']['date']
    print(f"Last commit to 'main':")
    print(f"  SHA: {last_commit_sha}")
    print(f"  Date: {last_commit_date}")

    # Fetch the open pull requests
    pulls_response = requests.get(pulls_url)
    pulls_response.raise_for_status()
    pulls_data = pulls_response.json()

    # Extract PR count and branch names
    pr_count = len(pulls_data)
    pr_branch_names = [pr['head']['ref'] for pr in pulls_data]

    print(f"\nOpen Pull Requests: {pr_count}")
    for branch_name in pr_branch_names:
      print(f"  Branch: {branch_name}")
  except requests.exceptions.RequestException as e:
    print(f"An error occurred: {e}")


def fetch_public_repos_with_prefix(org, prefix):
  base_url = f"https://api.github.com/orgs/{org}/repos"

  # List to store matching repo names
  matching_repos = []
  github_token = os.getenv('GH_TOKEN')
  try:
    # Paginated requests for all repos
    page = 1
    headers = {"Authorization": f"token {github_token}"}
    while True:
      response = requests.get(base_url, headers=headers, params={"page": page, "per_page": 50})
      response.raise_for_status()
      repos = response.json()

      # If there are no more repos, break the loop
      if not repos:
        break

      # Filter repos that start with the given prefix
      for repo in repos:
        if repo['name'].startswith(prefix) and not repo['archived']:
          print(repo['name'])
          matching_repos.append(repo['name'])

      page += 1

    print(f"Repositories starting with '{prefix}':")
    for repo_name in matching_repos:
      print(f"  - {repo_name}")

    return matching_repos
  except requests.exceptions.RequestException as e:
    print(f"An error occurred: {e}")
    return []


# Replace with your desired repo's owner and name
fetch_github_repo_data("ministryofjustice", "opg-digideps")
fetch_public_repos_with_prefix("ministryofjustice", "opg-")
