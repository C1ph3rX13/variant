package dynamic

import "variant/log"

func FileAnalyzer(path string) {
	e, err := Entropy(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("entropy: %v", e)

	md5, err := HashMD5(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("md5: %v", md5)

	sha1, err := HashSHA1(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("sh1: %v", sha1)

	sha256, err := HashSHA256(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("sha256: %v", sha256)

	sha512, err := HashSHA512(path)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("sha512: %v", sha512)
}
