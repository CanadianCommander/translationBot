import * as fs from "fs";

// take JSON on standard in, convert to TS, and write to the specified file.
async function main() {
    const args = process.argv.slice(2)
    if (args.length != 1) {
        console.error("Need 1 argument. <ts file to overwrite>")
        process.exit(1)
    }

    // parse json then re-export for pretty formatting.
    const json = JSON.parse(fs.readFileSync(process.stdin.fd).toString())
    fs.writeFileSync(
        args[0],
        "/* eslint-disable import/no-anonymous-default-export */ \n" +
        "export default " + JSON.stringify(json, null, 2) + '\n',
        {
            flag: "w"
        });
}

main()