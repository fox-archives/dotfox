use serde_json;
use std::fs;
use std::process::Command;
use std::string::String;

fn remove_whitespace(s: &mut String) {
	s.retain(|c| !c.is_whitespace());
}

fn main() -> serde_json::Result<()> {
	let output = Command::new("sh")
		.arg("-c")
		.arg("git remote get-url origin")
		.output()
		.expect("failed to execute process");

	let mut origin = String::from_utf8(output.stdout).unwrap();
	remove_whitespace(&mut origin);

	let jsonstr = r#"
	{
		"name": "foo"
	}
	"#;
	// let text = fs::read_to_string("package.json")?.parse()?;

	let json: serde_json::Value = serde_json::from_str(jsonstr)?;

	// if json.is_ok() {
	// 	let p: JsonValue = json.unwrap()
	// } else {
	// 	println!("could not parse")
	// }

	// println!("{}", text);

	println!("Hello, world!");
}
