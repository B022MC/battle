const { getDefaultConfig } = require('expo/metro-config');
const { withNativeWind } = require('nativewind/metro');
const { createProxyMiddleware } = require('http-proxy-middleware');

const defaultConfig = getDefaultConfig(__dirname);

const config = {
  ...defaultConfig, // 不改变Expo定制的metro配置
  server: {
    ...defaultConfig.server,
    enhanceMiddleware: (middleware) => {
      return (req, res, next) => {
        if (req.url.startsWith('/web')) {
          return createProxyMiddleware({
            target: 'http://118.24.135.171:8000',
            changeOrigin: true,
            pathRewrite: { '^/web': '/' },
          })(req, res, next);
        }
        return middleware(req, res, next);
      };
    },
  },
};

module.exports = withNativeWind(config, { input: './global.css', inlineRem: 16 });
