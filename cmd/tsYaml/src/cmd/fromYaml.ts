import * as fs from "fs";
import * as YAML from "yaml";

// take YAML on standard in, convert to TS, and write to the specified file.
async function main() {
    const args = process.argv.slice(2)
    if (args.length != 1) {
        console.error("Need 1 argument. <ts file to overwrite>")
        process.exit(1)
    }

    // parse json then re-export for pretty formatting.
    const yaml = YAML.parse(fs.readFileSync(process.stdin.fd).toString(), {version: "1.1"})
    fs.writeFileSync(
        args[0],
        "/* eslint-disable import/no-anonymous-default-export */ \n" +
        "export default " + JSON.stringify(yaml, null, 2) + '\n',
        {
            flag: "w"
        });
}

main()