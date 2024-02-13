const { program } = require('commander');
const { deployCommand } = require('./deploy');

console.log('run blinders cli');

program.addCommand(deployCommand);
program.parse();
