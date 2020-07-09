import fs from "fs";
import path from "path";
import toVfile from "to-vfile";
import remark from "remark";
import remarkLicense from "remark-license";
import remarkTableOfContents from "remark-toc";
import remarkUsage from "remark-usage";
import remarkTitle from "remark-title";
import remarkValidateLinks from "remark-validate-links";
import remarkWikiLink from "remark-wiki-link";

const filePath = path.join(__dirname, "../../readme.md");
const folderPath = path.dirname(filePath);
const file = fs.readFileSync(filePath, { encoding: "utf8" });
remark()
	// @ts-ignore
	.use(remarkTableOfContents, {
		maxDepth: 3,
		tight: true,
	})
	.use(remarkLicense, {
		name: "Edwin Kofler",
		license: "Apache-2.0",
		url: "https://edwinkofler.com",
	})
	.use(remarkTitle, {
		title: folderPath,
	})
	.use(remarkValidateLinks, {})
	// .use(remarkWikiLink, {})
	// .use(remarkUsage, {})
	.process(file, (err, output) => {
		if (err) console.error(err);

		console.info(output);
	});
