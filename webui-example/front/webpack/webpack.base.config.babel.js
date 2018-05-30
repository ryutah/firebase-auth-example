import path from "path";

export default {
  entry: path.join(__dirname, "../src/app.js"),
  output: {
    publicPath: "/",
    path: path.join(__dirname, "../dist"),
    filename: "bundle.js"
  },

  target: "web",
  resolve: {
    extensions: ["*", ".js", ".jsx"]
  },
  module: {
    rules: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules|firebaseui/,
        use: [
          {
            loader: "babel-loader"
          }
        ]
      },
      {
        test: /\.css/,
        use: [
          {
            loader: "style-loader"
          },
          {
            loader: "css-loader"
          }
        ]
      }
    ]
  }
};
