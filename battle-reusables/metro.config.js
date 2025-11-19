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
          console.log(`[Proxy] ${req.method} ${req.url} -> http://127.0.0.1:8000`);
          return createProxyMiddleware({
            target: 'http://127.0.0.1:8000',
            changeOrigin: true,
            pathRewrite: { '^/web': '/' },
            onProxyReq: (proxyReq, req, res) => {
              console.log(`[Proxy] Forwarding to: ${proxyReq.path}`);
            },
            onProxyRes: (proxyRes, req, res) => {
              console.log(`[Proxy] Response status: ${proxyRes.statusCode}`);
            },
            onError: (err, req, res) => {
              console.error('[Proxy] Error:', err.message);
              res.writeHead(500, { 'Content-Type': 'application/json' });
              res.end(JSON.stringify({ code: -1, msg: `Proxy error: ${err.message}` }));
            },
          })(req, res, next);
        }
        return middleware(req, res, next);
      };
    },
  },
};

module.exports = withNativeWind(config, { input: './global.css', inlineRem: 16 });
