import {
	assertEquals,
	assertArrayContains,
	assertThrows,
} from "https://deno.land/std/testing/asserts.ts";
import {
	getRepoUri,
	getEmail,
	getRepoInfo,
	getReadmeFilename,
	setPackageJsonFields,
} from "./util.ts";

// Deno.test({
// 	name: "getRepoUri",
// 	async fn(): Promise<void> {
// 		const uri = await getRepoUri();
// 		assertEquals(uri, "git@github.com:eankeen/globe");
// 	},
// });

// Deno.test({
// 	name: "getEmail",
// 	async fn(): Promise<void> {
// 		const email = await getEmail();
// 		assertEquals(email, "24364012+eankeen@users.noreply.github.com");
// 	},
// });

Deno.test({
	name: "getRepoInfo for bad protocol fails",
	fn(): void {
		assertThrows(() => getRepoInfo("http@github.com:eankeen/eankeen"));
	},
});

Deno.test({
	name: "getRepoInfo for ssh protocol works",
	fn(): void {
		const uri = "git@github.com:eankeen/project";
		const repo = getRepoInfo(uri);
		assertEquals(repo.host, "github");
		assertEquals(repo.name, "project");
		assertEquals(repo.owner, "eankeen");
		assertEquals(repo.protocol, "ssh");
		assertEquals(repo.vcs, "git");
	},
});

Deno.test({
	name: "getRepoInfo for https protocol works",
	fn(): void {
		const uri = "https://github.com/eankeen/project";
		const repo = getRepoInfo(uri);
		assertEquals(repo.host, "github");
		assertEquals(repo.name, "project");
		assertEquals(repo.owner, "eankeen");
		assertEquals(repo.protocol, "https");
		assertEquals(repo.vcs, "git");
	},
});

Deno.test({
	name: "getReadmeFilename",
	async fn(): Promise<void> {
		const readmeFile = await getReadmeFilename();
		assertEquals(readmeFile, "readme");
	},
});

Deno.test({
	name: "setPackageJsonFields does the correct thing",
	async fn(): Promise<void> {
		const uri = "https://github.com/eankeen/project";
		const repo = getRepoInfo(uri);

		const packageJson = await setPackageJsonFields({}, repo);
		assertEquals(
			packageJson.homepage,
			"https://github.com/eankeen/project#readme"
		);
		assertEquals(
			packageJson.bugs.url,
			"https://github.com/eankeen/project/issues"
		);
		// assertEquals(
		// packageJson.bugs.email,
		// "24364012+eankeen@users.noreply.github.com"
		// );
		assertEquals(packageJson.repository.type, "git");
		assertEquals(
			packageJson.repository.url,
			"git@github.com:eankeen/project.git"
		);
	},
});
