const { Command } = require('commander');
const { deploy } = require('./deploy');

const program = new Command();

program
	.name('blinders')
	.description('A CLI tool to manage Blinders project')
	.version('0.0.1');

program.addCommand(deploy);

program.parse();
