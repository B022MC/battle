module.exports = function (api) {
  api.cache(true);
  return {
    presets: [
      [
        'babel-preset-expo',
        {
          jsxImportSource: 'nativewind',
          unstable_transformImportMeta: true, // 添加 import.meta 转换支持
        },
      ],
      'nativewind/babel',
    ],
  };
};
