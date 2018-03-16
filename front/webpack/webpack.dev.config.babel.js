import config from "./webpack.base.config.babel"

config.devtool = "inline-source-map"
config.devServer = {
  inline: true,
  contentBase: "./public",
  port: "3000",
  host: "0.0.0.0",
}

export default config
