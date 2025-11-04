import { JSEncrypt } from "jsencrypt";

// 公钥
const publicKey = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1FEFPxWg1zYquQZJmqW4svXPLEjYPQ7PvtgqZWkQ0R5LnjnSjLdaNmpl1tcYA9lCiz21IEUK1ROfldnhVtS9ehK3zV2F0p5+6Vz/lU+kukfe6VVnwRXKGX2pBbQ3GU0m4Ih3yJ2KTTeMy6ZfZ0EC7agDGa9qDp1bfxoehfCz2HrTZl2WJkQOK+ily6uWt1zJDcUYffGD433eGEah8ISf3VOTdFmQzFiXxilwlnVTAeILGSS0/WhXrOs3Xqk5wgp41mMZxQ/uYIrQN4KIYZFlFosS1xuOQiTkrPBAVcZen8pfEjLv9yigD49H/QogkrONvQgfOYdruaN5sGvAODckwIDAQAB`;

// 公钥加密
export const encrypt = (val: string) => {
  const encryptor = new JSEncrypt();
  encryptor.setPublicKey(publicKey);
  const encrypted = encryptor.encrypt(val);
  if (!encrypted) throw new Error("RSA encryption failed");
  return encrypted;
};
