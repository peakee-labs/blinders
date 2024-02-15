const { Command } = require('commander');

const deploy = new Command();

deploy
	.name('deploy')
	.description('wrapped commands to deploy Blinders on AWS with Terraform');

deploy.command('echo').action(async () => {
	console.log('Hello deploy');
});

module.exports = { deploy };
