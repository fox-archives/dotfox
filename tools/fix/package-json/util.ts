import * as fs from "https://deno.land/std/fs/mod.ts";
export interface repo {
	vcs: "git";
	protocol: "ssh" | "https";
	host: string;
	owner: string;
	name: string;
}

export async function run(cmd: string[]): Promise<string> {
	const p = Deno.run({
		stdout: "piped",
		stderr: "piped",
		cmd,
	});

	// await p.status();
	const output = await p.output();
	// await p.close();

	const text = new TextDecoder("utf-8").decode(output);
	return text.trim();
}

export async function getRepoUri(): Promise<string> {
	return await run(["git", "remote", "get-url", "origin"]);
}

export async function getEmail(): Promise<string> {
	return await run(["git", "config", "--get", "user.email"]);
}

export function getRepoInfo(uri: string): repo {
	// @ts-ignore
	const repo: repo = {};
	repo.vcs = "git";

	// protocol
	{
		if (uri.startsWith("git@")) {
			repo.protocol = "ssh";
		} else if (uri.startsWith("https://")) {
			repo.protocol = "https";
		} else {
			throw new Error("miscellaneous protocol detected. not 'https'");
		}
	}

	// host
	{
		if (repo.protocol === "ssh") {
			const left = uri.indexOf("@") + 1;
			// not ':' because we don't want '.com'
			const right = uri.indexOf(".");
			repo.host = uri.slice(left, right);
		} else if (repo.protocol === "https") {
			const str = uri.slice("https://".length);

			const left = 0;
			const right = uri.indexOf("/");
			repo.host = str.slice(left, right);
		}
	}

	// owner
	{
		if (repo.protocol === "ssh") {
			const left = uri.indexOf(":") + 1;
			const right = uri.indexOf("/");
			repo.owner = uri.slice(left, right);
		} else if (repo.protocol === "https") {
			const oldStr = uri.slice("https://".length + 1);
			const str = oldStr.slice(oldStr.indexOf("/") + 1);

			const left = 0;
			const right = str.indexOf("/");
			repo.owner = str.slice(left, right);
		}
	}

	// name
	{
		if (repo.protocol === "ssh") {
			const left = uri.lastIndexOf("/") + 1;

			repo.name = uri.slice(left);
		} else if (repo.protocol === "https") {
			const oldStr = uri.slice("https://".length + 1);
			const str = oldStr.slice(oldStr.indexOf("/") + 1);

			const left = str.lastIndexOf("/") + 1;
			repo.name = str.slice(left);
		}
	}

	return repo;
}

export async function getReadmeFilename(): Promise<string | undefined> {
	const files = ["README.md", "readme.md", "README", "readme"];

	for (const file of files) {
		if (await fs.exists(file)) {
			return file;
		}
	}
	return;
}

export async function setPackageJsonFields(
	packageJson: Record<string, any>,
	repo: repo
): Promise<Record<string, any>> {
	const readmeFile = await getReadmeFilename();
	// const email = await getEmail();

	packageJson.homepage = `https://${repo.host}.com/${repo.owner}/${repo.name}#${readmeFile}`;

	packageJson.bugs = packageJson.bugs || {};
	packageJson.bugs.url = `https://${repo.host}.com/${repo.owner}/${repo.name}/issues`;
	// packageJson.bugs.email = email;

	// TODO: add license
	// packageJson.license

	packageJson.repository = packageJson.repository || {};
	packageJson.repository.type = repo.vcs;
	packageJson.repository.url = `git@${repo.host}.com:${repo.owner}/${repo.name}.git`;
	// TODO: add directory for monorepos
	// packageJson.repository.directory

	return packageJson;
}
