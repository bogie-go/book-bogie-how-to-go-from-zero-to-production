# Chapter 20: Work with an AI Assistant

This folder is not a full app. It holds only the two files Chapter 20 adds, a Claude Code hook that refuses hand edits to the files Bogie's generators own:

- [`.claude/hooks/guard-generated.sh`](.claude/hooks/guard-generated.sh): the guard script.
- [`.claude/settings.json`](.claude/settings.json): the project settings that run the guard before every `Edit` or `Write`.

Both are exactly the listings in the chapter. They go into your existing `depot` app, beside `AGENTS.md`. If you don't have one, use the finished shop in [Chapter 17](../ch17-keep-the-shop-healthy/). Its `AGENTS.md` and `CLAUDE.md` are the ones the chapter quotes. Bogie version: v1.0.0.

## Add the guard to your app

From the `depot` folder (it uses `jq`: `brew install jq` on macOS, or `sudo apt install jq` on Linux and inside WSL):

```sh
mkdir -p .claude/hooks
cp ../ch20-work-with-an-ai-assistant/.claude/hooks/guard-generated.sh .claude/hooks/
cp ../ch20-work-with-an-ai-assistant/.claude/settings.json .claude/
chmod +x .claude/hooks/guard-generated.sh
git add .claude && git commit -m "Guard the generators' files from hand edits"
```

Adjust the `../ch20-work-with-an-ai-assistant` path to wherever you cloned this repository.

## Test it

As in the chapter, from the `depot` folder:

```sh
export CLAUDE_PROJECT_DIR=$(pwd)
echo "{\"tool_name\":\"Edit\",\"tool_input\":{\"file_path\":\"$(pwd)/db/queries/products.sql\"}}" \
  | .claude/hooks/guard-generated.sh; echo "exit $?"
```

The guard refuses with exit 2 and says to use the `_extra` file, or `bogie g migration` and then `bogie model:sync [TABLE]`. The same check on `db/queries/products_extra.sql` exits 0 without printing anything.
