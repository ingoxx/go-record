package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	//encrypted, err := Encryption()
	//if err != nil {
	//	fmt.Println("加密错误:", err)
	//} else {
	//	fmt.Println("加密结果:", encrypted)
	//}

	encrypted := "YpbcPNR6Hc3WUU1BZPB9dnF44upJWzcCsOOh814QpjJw3UvkNzYMYRPBy2pwVMKYZ5UhCFt2LstEfPHOO8CIelNWy0NkJqMylgovKqsv/L4PuDbHnN1+/kChdNXoaxik2TqCvAz3Vitn4ckNzbgTQO1oA9M6tra3b80tbyJTKaK022jjjQsRLkso7shbYh72sLRWpu6pfyoU+xhllK2JMtUNBhYMVBInjq3Gb8RSfynjwhZf6E6IKRcVqxTDCJlqICi8/HMBXjdxCb/vxffeOk/ZkPWkTeMbFXDo6DAtENeLAjigMgXhYK+TSiSHDb187QkwfwlbfKH+v4B1tQJGP1RWntkSLQ63PkndYsFtumOIHRKjOfishPZs2Eu+QtFuxotdG5agmsFp6cFO/cVursdlXviwCKfgQ1q+mlrB8yAvv1hQKbOykSON1pKEuGLNreyA9cQiOF+pXKfTfpwVjtM9nCfgxJzLB9sgm1HIXqB60vmikJz1LOad1epWpEYeCo1Mj792gV2FsIaym7E4RsgkoPTZONUUSXip5bKxkaDIx9d1tEfUkVYd+7XH1GQrLWxlRraLPOcm7fBTY86wv0wJWgQnnenby4OpOgMtmwOgXfSCRnzfdx2g+Nt8xSQfj/JoT9gPGQ8ymWaVLW2jlQ6zboQSCBcrnt7GIGei3YKFkXFV0mmMCLpJSBsnwPtg9dOpdesdTkNtJMUYq0TBecexJjeasbHqHjERwzbHK7MJeWI8lAbMlxvPAINX6cslbuNwJscLfN9jOZD3yOgX50iWKs/50Lj2E7tJxJ1bqhztDYOoDPfrSlDavAUd5ttwl9SNFJ4U026sO36eSfY2xDM2lOIZCEOzZgx3u2DzpI/oKEpKQo318YQDprdTT/3ZTQruyzu59xmNuLgGxySz5mFpU2jqsWMolJtr7GdlZusPwkvAZlha5LjuaBFUaKsUBe55z1hBatww2cnAGcOGTMZXjz/QxVOYm/OSfKwN791Brkt0dsKYZ2uTlmfzKRO9RXS8VoSC8dM9TMEUO9xboN/Kdp/XRiiQasOAooUTPI3tATryJC5UdHrOTCKCakGvgoeWblmJ7FEVnRGpWALie1rV0Fqf7EHE733a04kZGzA27DcGtX3uc6nbXaehOHSxkv1RUzQ8QrAuRNLonJ+oVAAd2kPnuhEfNKS6OqPNTGoE8GyOaaCnhLXJYLVMseS+cPq4XboZBfmbyfECXVhB+t10Yc9w6DZYHdv1d3TxbN8B8G5+NVQOtP3e/cq/pjpMYd/wTAsfGgmc7CRQ9cHbwcqgN3vp0S9k8T29Zo4ny3JvUpeHTkYxZkX1qg4wiiXue3/6DG4aLaqMuUOHeAhlIaSJA/cFH61Kw9cW/wvaeH2xa86CwaPVP07TMDYumxBO/5vqqNGqMcuoIMreMjvPrOe4dWlKnAzuGvKqTDrEYus/YqMU/lF3zdIl9Pk3QF1TKP3qD0DdCNwsmc3c0ifhpaPFbvwmMo0b5DcUU6oXPnGL2HsK+VDeLQH4sliz+0btz42S7pOLEcFNzFhz+g8gkJKCsmqHXCY0Y4EzkCAWy/snBude3cIoq3st9AtYOQVrz0PL9b4Tvqlb08YQQPwfhQTopJJXzeyYCKpRacLGeeaYd0b+qrZTjyprFJacqVFbM97YYq2TCWWclIlIrtETj3Nuci6JMZuU9gl1Mu1R+XWSck5Fhcy05v1cUjYd9ZibWEQo+2FtYpMV0Cbyvfa1yvR08MgSjVBMqk0W7Oo7vAf+UGlAfx0zOxYLBK3F4Kdva4FSGIqTZAbGjmRaPMZBcLws0ohajTbn5xJClJuElD8DwI6TenKL7Q8GqjwhtFnYpJMkQt+TfI7TDPk7idIvhPOoV5Nvkp5PuObNHn/BYimMDXewCgJEoHBY1PdjrRYzrFXiaiWyAgAsm1BDvZKd/LRcajLKEKOb1nxbhxan6i+y9OQ3CY3IrnzdDtIdzaTFbxKrUvcH1PwIUiSG/4Fr6uPkFeEp5QekWlW+F2F+QWWOQQ3Bar7VV4jN7+rYwQnLaAq1ycjNXY1ABsjNflEpv3jI4CYd0R4xmBOqeMNzJLzzTR7l9rmkd6rME7Qb96/JkWF/RbSN0PqsXVlrlsXIiWXqT6EZS++CmQGYpEvH1hX80pKCrlSAeUjUuBNTu7M+eyAnqwPyHHDKLwpEzywJGxxURim0k/g7tIsqw612yt1kzjqMHo6NP73gySZWxNLfuhuTGITWZ4Lgm4iZAq9Om/Cf2WEteKqqBLn6RRTdslH0CcmzL5YYS1TX10ZqOAL0uTgGQznsFDDGxoyr7d2MdI6DQXqQaYbadhhAYO5xomRL0R56zDtkrH7k8KsOOuPpmtzCJS+GLt/Y6p6jNwyW3x8KJvembnefgeMRFgihT2QOnnXxOtSFMbllyLVfIj3y8i9dbmYsbzfBRR7QTs5eCWFagkUB0Lvi8seymmvGn+COKXyAFPmqJJM4X4Y97sH5gy6z0uwtkmcQJEzti34kdnIXasAdiMTk6P2e3CtRodyo4G8Vjj7+KdEbLJDF/OfFdRbqArr7+hoyJbqkHulR7H+2QXFkx0CtgmyjM4dzZo6YhIpXLPw7290Uy/VfULm96z4q1UrluZskCFRmLIcWISeT75RnpMZ+22AL7Ao5A8ysQkYFynq80Vuv0skFp0UsWLWsoEessjHR1HSnOUyntuxoSQsj+FxNOBI8XL3i+/cMRhLzRSvailUNMlI9l3dlHMKXq4bAYwzE6NpxPqUvw/q1DjOLHRzdoqTUZUquAVe5Y0ZiT8JOQbk0Cj8jyjltB9Xkdjkg8PkCf8zgmy30D6vhfEQOVBqcMIyTsADqW/5grfxvZZgaf5DAK/qa4u1gANbkUQYLSYTQjdPeyF2bTV7oVgW8XiASHw0FGy6NH11NJUVQKIA8oJi/p34Uz41WpdCWMJK1bSU7NRyf4V+697iC/p3+1/naUQ0J1Gt0Kw8TUUi0tLZLHmjNNP5vUCzG6uoqR9fbrJJiWITf62Uql1k5GdtQGDWm6Zot7dDsnpymGjmFCjJoR4kSZKJBWlDITlvPaD058zdDzK5N2RYbbMAXT4xA4sff5C6TI/efKVL57+c62CiQvXziViP3M4cSnLHzV0Rw0JJtRivqDUMFoUM+BUIsYYGUNQ9AZx7gxSezzOJjIeDSji28FAGpmrkZ/H4C98Y77Qx3nj418bM4LUIq5JQ2GsmcHsGb1x6EytNYRxBEt2W6v2ju6AWTd05piCFRVHYCdC6gqL9iagY2vQXZILg3ADREDAz5NDyUp4I5kyidyUdGZhGicK5qrxU0LyWgXbp9MvDNrjwb1fjvOELdt5ZrKxVNekb7EPnC8VLjBf3U92gWDXT9txWqTwS5olaJlMyE7miGwLuagoxb1bOeR+iRCK1F/SHz30QbqllcvMD+GN2yMuM5cqEkeic3/rG3TH7Tfk4ZDANjYoTm31AKp6Wmu251XVZ4BqbMGTuAnUiNmo4THUnOaZZkLLo/5Zs3gmGdMWb41Kb9VqsFZ0JRCoNWejfo0B6vuphQAcf+683yq+4C/3X1cDJ69rTTGPuFCKv0w5haGFjS+Mp5QOa5fnPi3ZBsLKKJW0G/wbUBwvrjlC4KnKx5AI97Tk4d/Wu4UrlGdXHQ0VClJmoDLKZe1os+ykI5ClUGoWk3mKy4MMRApWobk4RI1U9a8hJUXnikw8xjGxkTcP76ZUqbetcjXqNoLCPt21fimJ0ibdI/GOskuIoJr+l0/C7QxNMrsT2YV2pelAIYUtt+mCFNUj9Q3vwJTn3Hx+S6o/mxfkTiRHDigslr0n0WgURl+/bbRXk/ugRJ4wvPRthUIkWFOS+TnuR/3LUe75g8VkeRjQAkrcMXllwsDHafyHqIkCnGijnmHJ0F7nhSrLpaTM/8nT0+5VnJSUpZ00OCSixzUoxPwbWc22NaAqIP9nlBGyDbgA0IEB1ZEpBk/LxyECH8Yt94obViHIyZtjVTItiDucPmeN2D0wlULBVTl8B76X2IVdA4Ke1oALXUNSt7QL9+u1ZA/AIb4VeaXMlP8ij8an1y2rPe+Gz919bAbQtHkH2Kv7KeSJTDI9VR13uuLl0ObM8woFriiTuMDMbWcCKuYZN6A6dCQ5+GUhRQeVOV1wcSD7gSWxVed88u89lqhvaJXHNDwQ2+SeKUXnSm0+kOynXdJqbifxEevQbpaf6X1Ac/uvHoJzbfdjbD/w4rAHu5Wk+vytoX2FfMc+adGbtu4yC2Z67wXLvQGI7N9qP84VycZWeKqZ1J6fNF3rdjWONrxU79sfLDBlJTX6oof7CaNsN+ByGHJZ3N53FJos9m0TZmXAcQ9U7yI9isObZ7XEaKyBthTfnH3EOao42BV/PFZV4H70ftN56+F9/cg/MralX9xT18pSXDCdIIKLJFRZiUWsgPgzDmsceXIroh7CAOrtBn2t2F+PDAZO0Uy1w03YWOVkK1YQjd7CRzViJDidLB0izXvytx3pJwdzZrcVx957ymkVtkOkUsNbtSR6rRB0WlhX6zGUYcAv0KCWlaoZfIA0V1vECvW4TlMVRyh0+E3UeyC1g3Sm+ELvB9Wmkz1lirZu8DfQeEIG7ZhcvkkL9rmNZcDoZPIspTow85M8QINe3rSGv16ofMc5CHr/gjr4o60m0w4ANmf9Gjn1MtnWR7qgCaml6qkHMLTaNEoUMTjTX0Lh29VzPFIOml5+gRrskYGaDTeJ6NOZchRGs21jqg1kJX/G7G69H4SNMWHwiD9Rzir1Oc2U306S5Mhf4x94njhhi5FOSeWh1Uc4tsqpSFWrhEefi6J5g4C4rb2Hy/KRvHt4WkT9FTDCJH/5yiGv01gWP/9B2bK4Uzy83YuTbYGk/q+SsNfpIHwC3UCMFXzvjJ7yBaSnJX3YW1lynYjGHjxb4O7f4GFXC6+3yLYQKczNnzZNWzrxlQkp7WjJIbgqBlGL9i3jpdfcj613gUrYe6yvczg/is1tCnbZ1AhkEDNuAklFl88Ge72XIxoClrDjJu5sHXG/uC3IaufhVGOlHTggsZ5n9uqPWdMUqlLkVRCZfClHfMLG2vV3Aer/2H"

	key := []byte("my-secret-key-32bytesaaaaaa-12313123!!")
	key = key[:32]
	decryptAES, err := DecryptAES(encrypted, key)
	if err != nil {
		fmt.Println("解密错误:", err)
	} else {
		fmt.Println("解密结果:", decryptAES)
	}
}

func Encryption() (string, error) {
	// **确保密钥是 16、24 或 32 字节**
	key := []byte("my-secret-key-32bytesaaaaaa-12313123!!") // 37 字节，超长了
	key = key[:32]                                          // 取前 32 字节（AES-256）

	// 需要加密的数据
	plaintext := "aabbc"

	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES 失败: %w", err)
	}

	// 生成 IV（初始化向量），长度必须与 block.BlockSize() 相同
	iv := key[:block.BlockSize()]

	plainBytes := []byte(plaintext)
	ciphertext := make([]byte, len(plainBytes))

	// 使用 CTR 模式加密
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plainBytes)

	// 返回 Base64 编码的加密结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// **AES 解密（CTR 模式）**
func DecryptAES(cipherTextBase64 string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("创建 AES 失败: %w", err)
	}

	// 生成 IV（初始化向量），长度必须与 block.BlockSize() 一致
	iv := key[:block.BlockSize()]

	// 解码 Base64
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %w", err)
	}

	plainText := make([]byte, len(cipherText))

	// 创建 CTR 解密流（CTR 其实是对称的，只需再次 XOR 即可解密）
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
