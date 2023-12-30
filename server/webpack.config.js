const path = require("path");
const glob = require("glob");

module.exports = {
  mode: "production",
  //   entry: { lib: glob.sync("./public/src/lib/*.js") },
  //   output: {
  //     path: path.join(__dirname, "/public/dist/lib"),
  //     filename: "[name].bundle.js",
  //     sourceMapFilename: "[name].map",
  //   },
  entry: glob.sync("./public/src/lib/*.js").reduce(function (obj, el) {
    obj[path.parse(el).name] = el;
    return obj;
  }, {}),
  output: {
    path: path.join(__dirname, "/public/dist/js"),
    filename: "[name].js",
  },
};
