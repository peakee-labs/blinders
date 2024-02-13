const { Command } = require('commander');

const deployCommand = new Command('deploy').help('');

module.exports = { deployCommand };
