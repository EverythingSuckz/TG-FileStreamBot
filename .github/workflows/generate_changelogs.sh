#!/bin/bash
# set -x

echo "## CHANGELOGS"

last_commit=$(curl -s https://api.github.com/repos/EverythingSuckz/TG-FileStreamBot/tags | jq -r '.[0].commit.sha')

resp=$(curl -s https://api.github.com/repos/EverythingSuckz/TG-FileStreamBot/commits?per_page=5)
data=$(echo "${resp}" | jq -c '.[]')
while read -r i; do
    if [[ $i == $last_commit ]]; then
        break
    fi
    sha=$(echo $i | jq -r '.sha')
    message=$(echo $i | jq -r '.commit.message' |  tr '\n' ' ')
    commit_url=$(echo $i | jq -r '.html_url')
    author_name=$(echo $i | jq -r '.commit.author.name')
    author_url=$(echo $i | jq -r '.author.html_url')
    echo "- [$(echo $sha | cut -c1-7)]($commit_url) - $message by [$author_name]($author_url)"
done <<< "$data"

if [[ ! -z "$last_commit" ]]; then
    last_commit="main"
fi

compare_url="https://github.com/EverythingSuckz/TG-FileStreamBot/compare/$(echo $resp | jq -r '.[0].sha')...$last_commit"
echo -e "- [.....]($compare_url)\n\nView Full Changelogs [here]($compare_url)"
