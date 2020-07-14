// @ts-ignore
import fs from "fs";
// @ts-ignore
import remark from "remark";
// @ts-ignore
import visit from "unist-util-visit";
// @ts-expect-error
import vfile from "vfile";

const filepath = "README.md";
const content = fs.readFileSync("../../README.md", "utf8");

const owner = "eankeen";
const repo = "fox-suite";

remark()
	// .use(checkTitle)
	.use(function plugin(options) {
		return function transform(ast, file) {
			let doNext = false;
			for (let node of ast.children) {
				if (doNext) {
					const n = node.children;
					if (!n) throw new Error("n is not truthy");
					node.children = [
						{
							type: "link",
							url: `https://starchart.cc/${owner}/${repo}`,
							children: [
								{
									type: "image",
									url: `https://starchart.cc/${owner}/${repo}.svg`,
									alt: "Stargazers over time",
								},
							],
						},
					];
					console.info(node);
					break;
				}
				if (
					node.type === "heading" &&
					node.children?.[0]?.value === "Star Chart"
				) {
					console.info(node);
					doNext = true;
				}
			}
		};
	})
	.process(vfile({ file: filepath, contents: content }))
	.then(async (vf) => {
		console.info(vf);
		fs.writeFileSync("../../README.md", vf.contents);
	})
	.catch((err) => {
		console.error(err);
		process.exitCode = 1;
	});
