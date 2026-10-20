#!/usr/bin/env node
const path = require("path");
const { spawnSync } = require("child_process");

const cliPath = path.join(__dirname, "../src/cli.ts");
const tsxBin = path.join(__dirname, "../node_modules/.bin/tsx");

const result = spawnSync(tsxBin, [cliPath, ...process.argv.slice(2)], {
  stdio: "inherit",
});

process.exit(result.status ?? 0);
