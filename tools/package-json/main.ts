import * as fs from "https://deno.land/std/fs/mod.ts";
import { getRepoUri, getRepoInfo, setPackageJsonFields, repo } from "./util.ts";

const uri = await getRepoUri();
let repo: repo;
try {
	repo = getRepoInfo(uri);
} catch (err) {
	console.error(err);
}

let packageJson: any = await fs.readJson("package.json");

packageJson = await setPackageJsonFields(packageJson, repo);
console.log("ffooo", packageJson);
await fs.writeJson("package.json", packageJson, {
	// TODO: don't assume indentation type
	spaces: "\t",
});
