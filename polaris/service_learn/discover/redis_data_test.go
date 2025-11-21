package discover

import (
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	v2 "git.woa.com/polaris/polaris-server-api/api/v2/common"
	"github.com/golang/protobuf/proto"
	"testing"
)

func TestRedisData(t *testing.T) {
	pbSourceStr := "\n\xf1\x04\n\x0f192000065:83357\x12\nProduction\x1a#\n\x10internal-cl5-sid\x12\x0f192000065:83357\x1a\x1e\n\x16internal-enable-nearby\x12\x04true\x1a\x14\n\x0cinternal-cl5\x12\x04true\x1a\x1c\n\x11internal-cl5-bid4\x12\a1140601\"\xb7\x0210074;10274;10697;10693;10045;10049;10333;10096;10140;10020;10017;10023;10048;10088;10845;10080;11320;10537;10694;10294;10615;10704;10290;10241;10466;10201;10242;11631;10225;10162;10257;10213;10406;10170;10490;10267;10663;10175;10204;10745;10868;10165;11150;10513;10072;10276;10117;10569;10529;10465;10277;10394*\x19UserData.main_port.\xe6\xb7\xb1\xe5\x9c\xb32\x16PCG-\xe5\xa4\xa7\xe6\x95\xb0\xe6\x8d\xae\xe5\xb9\xb3\xe5\x8f\xb0\xe9\x83\xa8Z+jimpang;barrysun;boshao;zhizhuoliu;flashguo`\xde\xb4\xdc\x81\x06h\xaf\x82\xf5\xc8\x06r c599f3f1aed741048bc562a6c02e6dffz\x11123_plugin_naming\x12\xd0\x02\n(84a12be606b0624a2e7dcec06ef69e2998c2e1bc\x1a\x0e11.177.147.138 \xeb\x86\x01@dX\x01j-\n\x06\xe5\x8d\x8e\xe5\x8d\x97\x12\x06\xe6\xb7\xb1\xe5\x9c\xb3\x1a\x14\xe6\xb7\xb1\xe5\x9c\xb3\xe7\xa7\xbb\xe5\x8a\xa8\xe8\x8d\x94\xe6\x99\xafDC \x02(\x1e0\xa4gr\x1b\n\x12internal-cl5-setId\x12\x05NOSETr\r\n\aip-type\x12\x02v4r/\n\x0econtainer_name\x12\x1dformal.IWAN.UserData.sz100032r\r\n\x03env\x12\x06formalr@\n\rinstance_name\x12/cls-9zmchzdc-4b5e3c47797bff4c19b0efa8cd9f01fd-0\x80\x01\x97\xf3\xde\xbc\x06\x88\x01\xab\x85\x91\xc5\x06\x92\x01 896620dca6654440b70d2a9a782ef800\x12\xce\x02\n(e1e27cb791eb2873fdb4bdd3b2d80203dd3c5ed6\x1a\r9.146.133.126 \xdas@dX\x01j-\n\x06\xe5\x8d\x8e\xe5\x8d\x97\x12\x06\xe6\xb7\xb1\xe5\x9c\xb3\x1a\x14\xe6\xb7\xb1\xe5\x9c\xb3\xe7\xa7\xbb\xe5\x8a\xa8\xe5\x85\x89\xe6\x98\x8eDC \x02(\x1e0\xc5Hr@\n\rinstance_name\x12/cls-560o61y4-f1ee034db57a7dd6ed1187c1aedb5f4c-0r\x1b\n\x12internal-cl5-setId\x12\x05NOSETr\r\n\aip-type\x12\x02v4r/\n\x0econtainer_name\x12\x1dformal.IWAN.UserData.sz100038r\r\n\x03env\x12\x06formal\x80\x01\xa0\xce\xd8\xc1\x06\x88\x01\xab\x85\x91\xc5\x06\x92\x01 c5a5930be8eb406b94a480d8aba4919e\x12\xc4\x02\n(54af81b5f217eb8250d9248b8afb4552638222ee\x1a\x0e11.160.134.193 \xc5q@dX\x01j\"\n\x06\xe5\x8d\x8e\xe5\x8d\x97\x12\x06\xe6\xb7\xb1\xe5\x9c\xb3\x1a\x0c\xe6\xb7\xb1\xe5\x9c\xb3\xe5\x9b\x9b\xe5\x8c\xba \x02(\x1er@\n\rinstance_name\x12/cls-9zmchzdc-4b5e3c47797bff4c19b0efa8cd9f01fd-1r\x1b\n\x12internal-cl5-setId\x12\x05NOSETr\r\n\aip-type\x12\x02v4r/\n\x0econtainer_name\x12\x1dformal.IWAN.UserData.sz100035r\r\n\x03env\x12\x06formal\x80\x01\xa8\xf2\xcf\xbe\x06\x88\x01\xab\x85\x91\xc5\x06\x92\x01 c7acbc761ac445648b2651e475556bf7"
	svcIns := &v2.DiscoverResponse{}

	if err := proto.Unmarshal([]byte(pbSourceStr), svcIns); err != nil {
		t.Fatal(err)
	}

	t.Log(svcIns)
}

func TestRedisData1(t *testing.T) {
	pbSourceStr := "\n\xf1\x04\n\x0f192000065:83357\x12\nProduction\x1a\x14\n\x0cinternal-cl5\x12\x04true\x1a\x1c\n\x11internal-cl5-bid4\x12\a1140601\x1a#\n\x10internal-cl5-sid\x12\x0f192000065:83357\x1a\x1e\n\x16internal-enable-nearby\x12\x04true\"\xb7\x0210074;10274;10697;10693;10045;10049;10333;10096;10140;10020;10017;10023;10048;10088;10845;10080;11320;10537;10694;10294;10615;10704;10290;10241;10466;10201;10242;11631;10225;10162;10257;10213;10406;10170;10490;10267;10663;10175;10204;10745;10868;10165;11150;10513;10072;10276;10117;10569;10529;10465;10277;10394*\x19UserData.main_port.\xe6\xb7\xb1\xe5\x9c\xb32\x16PCG-\xe5\xa4\xa7\xe6\x95\xb0\xe6\x8d\xae\xe5\xb9\xb3\xe5\x8f\xb0\xe9\x83\xa8Z+jimpang;barrysun;boshao;zhizhuoliu;flashguo`\xde\xb4\xdc\x81\x06h\xe9\xc1\xf5\xc8\x06r c599f3f1aed741048bc562a6c02e6dffz\x11123_plugin_naming\x12\xc4\x02\n(54af81b5f217eb8250d9248b8afb4552638222ee\x1a\x0e11.160.134.193 \xc5q@dX\x01j\"\n\x06\xe5\x8d\x8e\xe5\x8d\x97\x12\x06\xe6\xb7\xb1\xe5\x9c\xb3\x1a\x0c\xe6\xb7\xb1\xe5\x9c\xb3\xe5\x9b\x9b\xe5\x8c\xba \x02(\x1er@\n\rinstance_name\x12/cls-9zmchzdc-4b5e3c47797bff4c19b0efa8cd9f01fd-1r\x1b\n\x12internal-cl5-setId\x12\x05NOSETr\r\n\aip-type\x12\x02v4r/\n\x0econtainer_name\x12\x1dformal.IWAN.UserData.sz100035r\r\n\x03env\x12\x06formal\x80\x01\xa8\xf2\xcf\xbe\x06\x88\x01\xe2\xc1\xf5\xc8\x06\x92\x01 c7acbc761ac445648b2651e475556bf7\x12\xd0\x02\n(84a12be606b0624a2e7dcec06ef69e2998c2e1bc\x1a\x0e11.177.147.138 \xeb\x86\x01@dX\x01j-\n\x06\xe5\x8d\x8e\xe5\x8d\x97\x12\x06\xe6\xb7\xb1\xe5\x9c\xb3\x1a\x14\xe6\xb7\xb1\xe5\x9c\xb3\xe7\xa7\xbb\xe5\x8a\xa8\xe8\x8d\x94\xe6\x99\xafDC \x02(\x1e0\xa4gr\r\n\x03env\x12\x06formalr@\n\rinstance_name\x12/cls-9zmchzdc-4b5e3c47797bff4c19b0efa8cd9f01fd-0r\x1b\n\x12internal-cl5-setId\x12\x05NOSETr\r\n\aip-type\x12\x02v4r/\n\x0econtainer_name\x12\x1dformal.IWAN.UserData.sz100032\x80\x01\x97\xf3\xde\xbc\x06\x88\x01\xe2\xc1\xf5\xc8\x06\x92\x01 896620dca6654440b70d2a9a782ef800\x12\xce\x02\n(e1e27cb791eb2873fdb4bdd3b2d80203dd3c5ed6\x1a\r9.146.133.126 \xdas@dX\x01j-\n\x06\xe5\x8d\x8e\xe5\x8d\x97\x12\x06\xe6\xb7\xb1\xe5\x9c\xb3\x1a\x14\xe6\xb7\xb1\xe5\x9c\xb3\xe7\xa7\xbb\xe5\x8a\xa8\xe5\x85\x89\xe6\x98\x8eDC \x02(\x1e0\xc5Hr@\n\rinstance_name\x12/cls-560o61y4-f1ee034db57a7dd6ed1187c1aedb5f4c-0r\x1b\n\x12internal-cl5-setId\x12\x05NOSETr\r\n\aip-type\x12\x02v4r/\n\x0econtainer_name\x12\x1dformal.IWAN.UserData.sz100038r\r\n\x03env\x12\x06formal\x80\x01\xa0\xce\xd8\xc1\x06\x88\x01\xe2\xc1\xf5\xc8\x06\x92\x01 c5a5930be8eb406b94a480d8aba4919e\x12\x8f\x01\n(0065562e2629e935c70176becd1b70a282657207\x1a\t127.1.1.1 \xb8E*\x04grpc2\x02v1@dX\x01`\x01r\x14\n\atestkey\x12\ttestvalue\x80\x01\xd2\xc1\xf5\xc8\x06\x88\x01\xe2\xc1\xf5\xc8\x06\x92\x01 881324183ced4746a9c72df40048cc0c"
	svcIns := &v2.DiscoverResponse{}

	if err := proto.Unmarshal([]byte(pbSourceStr), svcIns); err != nil {
		t.Fatal(err)
	}

	t.Log(svcIns)
}

func TestRedisData2(t *testing.T) {
	pbSourceStr := "\n6\n4trpc.RondaCommonKVTemplateGZ.CacheServer1.RootClient\x12\r\n\x0bDevelopment\"\x85\x01\nX\n6\n4trpc.RondaCommonKVTemplateGZ.CacheServer1.RootClient\x12\r\n\x0bDevelopment\x1a\x0f\n\x03env\x12\b\x12\x06\n\x04test\x12)\n\x03\n\x01*\x12\r\n\x0bDevelopment\x1a\x0f\n\x03env\x12\b\x12\x06\n\x04test*\x02\bd*\x15\n\x132021-12-29 16:33:292\x15\n\x132025-11-13 08:27:14:\"\n b35e5a008b5b4d8081b1b9ce69974327"
	svcIns := &apiV1Model.Routing{}

	if err := proto.Unmarshal([]byte(pbSourceStr), svcIns); err != nil {
		t.Fatal(err)
	}
	t.Log(svcIns)
}
