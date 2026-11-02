package main

import "context"

type seedVillage struct {
	id, name, region, summary, history, cover string
	lat, lng                                  float64
	honors, themes                            []string
	sort                                      int
}

type seedSpot struct {
	id, villageID, name, category, summary, image string
	lat, lng                                      float64
	sort                                          int
}

func seedVillages() []seedVillage {
	return []seedVillage{
		{"village-shaoshan", "韶山村", "韶山市韶山乡",
			"毛泽东同志故里，全国红色旅游融合发展试点核心村域，以故居、铜像广场与纪念馆构成经典参观序列。",
			"韶山村是毛泽东同志的故乡。1925 年他在此开展农民运动，1927 年考察湖南农民运动时回到韶山。村内保存毛泽东同志故居、毛氏宗祠、南岸私塾等旧址，周边建有毛泽东铜像广场、韶山毛泽东同志纪念馆与烈士陵园，形成完整的革命纪念体系与文化教育场景。",
			"/static/red/village-shaoshan.png", 27.9146, 112.5283,
			[]string{"全国红色旅游融合发展试点", "全国乡村旅游重点村"}, []string{"伟人故里", "革命纪念", "研学教育"}, 1},
		{"village-wushi", "乌石村", "湘潭县乌石镇",
			"彭德怀同志故里，乌石峰下的红色村落，以纪念馆、故居与国防教育为主题特色。",
			"乌石村是彭德怀同志的故乡。村内的彭德怀纪念馆依山而建，与彭德怀故居、铜像广场、德怀墓园共同组成纪念序列。抗日战争与解放战争时期，彭德怀同志屡建战功，其“横刀立马”的形象成为革命英雄主义教育的经典素材，村域内保留了大量与彭总成长相关的场景与口述记忆。",
			"/static/red/village-wushi.png", 27.7361, 112.8147,
			[]string{"国家国防教育示范基地", "省级红色旅游经典景区"}, []string{"元帅故里", "国防教育", "英雄故事"}, 2},
		{"village-quanhu", "泉湖村", "湘乡市龙洞镇",
			"陈赓同志故里，湘中丘陵间的红色村落，以故居与将军事迹展陈为核心内容。",
			"泉湖村是中国人民解放军大将陈赓同志的故乡。陈赓同志历经北伐、长征、抗日战争、解放战争与抗美援朝，是新中国国防科技教育事业的重要奠基人之一。村内陈赓故居保存较为完整，周边保留湘中传统民居风貌与农耕景观，适合开展将军事迹研学和乡村体验活动。",
			"/static/red/village-quanhu.png", 27.7861, 112.4189,
			[]string{"省级文物保护单位", "爱国主义教育基地"}, []string{"大将故里", "军工记忆", "乡村研学"}, 3},
		{"village-xingtang", "杏花村", "湘潭县白石镇",
			"齐白石故里，星斗塘与借山吟馆所在地，是湘潭独有的艺术大师文化标识。",
			"杏花村是艺术大师齐白石先生的诞生地。白石老人从木匠成长为一代画坛宗师，其“衰年变法”与“不教一日闲过”的治艺精神，成为湘潭最具辨识度的艺术文化资源。村内保留星斗塘、白石铺等与大师成长相关的场所，周边以湘莲种植与乡村田园景观为特色。",
			"/static/red/village-xingtang.png", 27.6536, 112.9517,
			[]string{"艺术大师故里", "省级乡村旅游重点村"}, []string{"大师故里", "艺术研学", "湘莲田园"}, 4},
		{"village-yinshan", "银田村", "韶山市银田镇",
			"韶山灌区与银田老街交汇的红色村域，以水利建设记忆与乡村振兴实践为特色。",
			"银田村地处韶山灌区沿线，二十世纪六十年代修建的韶山灌区是新中国水利建设的重要工程，被誉为“南方的红旗渠”。村域内银田老街保留湘中集镇格局，通过红色水利记忆与乡村产业更新相结合，形成了以研学、农事体验和农产品加工为支撑的发展路径。",
			"/static/red/village-yinshan.png", 27.8625, 112.5936,
			[]string{"乡村振兴示范村", "水利工程遗产地"}, []string{"水利记忆", "老街风貌", "助农产业"}, 5},
		{"village-dongshan", "东山社区", "湘乡市东山街道",
			"毛泽东同志少年求学地东山学校所在片区，以求学立志与教育振兴为主题。",
			"东山学校创建于 1895 年，是毛泽东同志少年时期离开家乡求学的地方。1910 年他进入东山高等小学堂读书，接触到新学与进步思想，立下“改造中国与世界”的志向。校园内保留学堂旧址与历史建筑，是开展青少年立志教育与研学活动的重要场所。",
			"/static/red/village-dongshan.png", 27.7348, 112.5264,
			[]string{"全国重点文物保护单位", "青少年励志教育基地"}, []string{"求学立志", "新学启蒙", "校园研学"}, 6},
	}
}

func seedSpots() []seedSpot {
	return []seedSpot{
		{"spot-maoguju", "village-shaoshan", "毛泽东同志故居", "革命旧址",
			"湘中传统农舍建筑，陈列少年毛泽东生活场景与家庭史料。", "/static/red/spot-maoguju.png", 27.9128, 112.5269, 1},
		{"spot-maotongxiang", "village-shaoshan", "毛泽东铜像广场", "纪念场所",
			"广场以铜像为主体，是开展献花、宣誓与现场教学的主要场地。", "/static/red/spot-tongxiang.png", 27.9161, 112.5310, 2},
		{"spot-maozhendang", "village-shaoshan", "韶山毛泽东同志纪念馆", "专题展馆",
			"系统展示毛泽东同志生平与革命实践，兼具文物陈列与教育功能。", "/static/red/spot-museum.png", 27.9175, 112.5323, 3},
		{"spot-pengjinianguan", "village-wushi", "彭德怀纪念馆", "专题展馆",
			"依山而建的纪念建筑群，展示彭德怀同志生平与军事贡献。", "/static/red/spot-pengjinianguan.png", 27.7350, 112.8124, 1},
		{"spot-pengguju", "village-wushi", "彭德怀故居", "革命旧址",
			"保存彭德怀同志少年生活场景，周边保留乌石峰自然景观。", "/static/red/spot-pengguju.png", 27.7368, 112.8159, 2},
		{"spot-defengmu", "village-wushi", "德怀墓园", "纪念场所",
			"庄重肃穆的纪念空间，适合开展敬献花篮与主题教育活动。", "/static/red/spot-defengmu.png", 27.7341, 112.8137, 3},
		{"spot-chengengguju", "village-quanhu", "陈赓故居", "革命旧址",
			"湘中民居格局的将军故里，陈列陈赓同志生平事迹。", "/static/red/spot-chengengguju.png", 27.7868, 112.4175, 1},
		{"spot-quanhu-tian", "village-quanhu", "泉湖农耕体验区", "乡村体验",
			"稻田与菜园组成的农事体验区，用于研学劳动实践。", "/static/red/spot-quanhu-tian.png", 27.7849, 112.4202, 2},
		{"spot-xingdoutang", "village-xingtang", "星斗塘与齐白石故居", "艺术旧址",
			"齐白石诞生地与成长环境，保留湘中农家院落与荷塘景观。", "/static/red/spot-xingdoutang.png", 27.6528, 112.9503, 1},
		{"spot-jieshan", "village-xingtang", "借山吟馆艺术陈列", "艺术展馆",
			"展示齐白石书画复制品、篆刻与从艺历程。", "/static/red/spot-jieshan.png", 27.6549, 112.9531, 2},
		{"spot-yintianlaojie", "village-yinshan", "银田老街", "历史街区",
			"湘中集镇格局保存区，串联老街商铺与红色水利记忆。", "/static/red/spot-yintianlaojie.png", 27.8617, 112.5921, 1},
		{"spot-guanqu", "village-yinshan", "韶山灌区纪念点", "水利遗产",
			"展示灌区建设历程与劳动精神，适合开展主题研学。", "/static/red/spot-guanqu.png", 27.8643, 112.5958, 2},
		{"spot-dongshantang", "village-dongshan", "东山学校旧址", "求学旧址",
			"毛泽东同志少年求学地，保存学堂旧貌与励志故事。", "/static/red/spot-dongshantang.png", 27.7342, 112.5251, 1},
		{"spot-dongshan-jinianguan", "village-dongshan", "东山学校校史馆", "专题展馆",
			"展陈学校沿革与求学立志故事，配套研学课程。", "/static/red/spot-dongshan-museum.png", 27.7356, 112.5278, 2},
	}
}

type seedRoute struct {
	id, name, theme, summary, cover string
	days                            int
	distance                        float64
	spotIDs, spotNames              []string
	sort                            int
}

type seedStory struct {
	id, villageID, title, category, era, summary, content, cover, arImage string
	views                                                                 int
}

type seedProduct struct {
	id, name, category, cover, summary, detail, storyCard, farmerID, farmerName, villageID string
	price, original, stock, sales                                                          int
	tags                                                                                   []string
	customizable                                                                           bool
	sort                                                                                   int
}

type seedFarmer struct {
	id, name, villageID, villageName, intro, avatar, project string
	target, count                                            int
}

type seedTask struct {
	id, villageID, spotID, title, kind, summary, cover string
	points                                             int
	sort                                               int
}

func defaultPointsPrice(priceCents int) int {
	value := priceCents / 10
	if value < 180 {
		value = 180
	}
	return value
}

func defaultRewardPoints(priceCents int) int {
	value := priceCents / 500
	if value < 5 {
		value = 5
	}
	return value
}

func defaultRewardExp(priceCents int) int {
	value := priceCents / 400
	if value < 10 {
		value = 10
	}
	return value
}

func seedRoutes() []seedRoute {
	return []seedRoute{
		{"route-shaoshan-classic", "韶山经典红迹线", "伟人故里",
			"串联故居、铜像广场与纪念馆，一日完成从生活场景到生平叙事的完整参观。",
			"/static/red/route-shaoshan.png", 1, 8.5,
			[]string{"spot-maoguju", "spot-maotongxiang", "spot-maozhendang"},
			[]string{"毛泽东同志故居", "毛泽东铜像广场", "韶山毛泽东同志纪念馆"}, 1},
		{"route-wushi-hero", "乌石元帅足迹线", "英雄故事",
			"以乌石峰为背景，串联纪念馆、故居与墓园，突出国防教育与英雄精神。",
			"/static/red/route-wushi.png", 1, 6.2,
			[]string{"spot-pengjinianguan", "spot-pengguju", "spot-defengmu"},
			[]string{"彭德怀纪念馆", "彭德怀故居", "德怀墓园"}, 2},
		{"route-art-village", "大师故里艺术研学线", "艺术乡建",
			"从星斗塘到借山吟馆，结合湘莲田园与写生实践，做一条艺术主题研学线路。",
			"/static/red/route-art.png", 2, 15.4,
			[]string{"spot-xingdoutang", "spot-jieshan", "spot-yintianlaojie"},
			[]string{"星斗塘与齐白石故居", "借山吟馆艺术陈列", "银田老街"}, 3},
	}
}

func seedStories() []seedStory {
	return []seedStory{
		{"story-qiushi", "village-dongshan", "走出韶山冲的少年", "求学立志", "1910 年",
			"十六岁的毛泽东离开韶山，步行来到东山高等小学堂，第一次系统接触新学。",
			"一九一〇年秋天，十六岁的毛泽东挑着行李离开韶山冲，来到湘乡东山高等小学堂求学。在此之前，他读过六年私塾，也在家里务农记账。进入东山学堂后，他第一次系统接触到地理、历史与自然科学，读到了《新民丛报》等进步刊物，并写下“孩儿立志出乡关，学不成名誓不还”的诗句。这段求学经历，是他从乡村少年走向革命者的重要起点。",
			"/static/red/story-qiushi.png", "/static/red/ar-dongshan.png", 1286},
		{"story-nongyun", "village-shaoshan", "韶山冲里的农民夜校", "革命事件", "1925 年",
			"毛泽东同志在韶山开展农民运动，办夜校、建支部，点燃湘中农运火种。",
			"一九二五年，毛泽东同志回到韶山养病期间，深入农户调查，先后在毛氏宗祠、李家祠堂等地办起农民夜校，向农民讲解为什么受苦、怎样才能翻身。这一年，中共韶山特别支部成立，成为湖南最早的农村党支部之一。短短几个月，韶山一带的农民协会迅速建立，为后来的湖南农民运动积累了经验。",
			"/static/red/story-nongyun.png", "/static/red/ar-shaoshan.png", 1043},
		{"story-hengdao", "village-wushi", "横刀立马彭大将军", "人物故事", "1947 年",
			"西北战场上的彭德怀同志，以“横刀立马”的形象为人们所铭记。",
			"一九四七年，国民党军重点进攻陕甘宁边区。彭德怀同志临危受命，指挥西北野战军以不足三万的兵力对抗十倍之敌，取得青化砭、羊马河、蟠龙三战三捷，稳定了西北战局。毛泽东同志曾写下“谁敢横刀立马，唯我彭大将军”的诗句。乌石村也因这位元帅而被更多人记住。",
			"/static/red/story-hengdao.png", "/static/red/ar-wushi.png", 967},
		{"story-baishi", "village-xingtang", "从木匠到画坛宗师", "艺术人生", "1864-1957 年",
			"齐白石从星斗塘的农家少年成长为一代画坛宗师，一生不教一日闲过。",
			"齐白石出生于湘潭县白石铺杏子塘一户农家，少年时学木匠，做雕花手艺。二十七岁起正式学画，五十五岁定居北京后实行“衰年变法”，自创红花墨叶一派，成为二十世纪中国画的重要代表。他常说“不教一日闲过”，一生作画数以万计，虾、蟹、花鸟与篆刻尤为世人所重。故里星斗塘至今保留着他成长的田园环境。",
			"/static/red/story-baishi.png", "/static/red/ar-xingtang.png", 1524},
		{"story-guanqu", "village-yinshan", "韶山灌区的建设记忆", "建设故事", "1965-1966 年",
			"十万建设者劈山引水，建成被誉为“南方红旗渠”的韶山灌区。",
			"二十世纪六十年代中期，为改变湘中地区干旱缺水的面貌，十万余名建设者参与韶山灌区工程。工程劈开山岭、架设渡槽，建成干渠与众多支渠，灌溉湘潭、湘乡、韶山等地大片农田。建设过程中形成了自力更生、艰苦奋斗的水利精神，银田一带至今保留着与灌区相关的工程遗迹与集体记忆。",
			"/static/red/story-guanqu.png", "/static/red/ar-yinshan.png", 736},
		{"story-xianglian", "village-xingtang", "一粒湘莲的百年光景", "风物故事", "当代",
			"湘莲是湘潭传统物产，也是当地乡村产业与餐桌记忆的重要组成。",
			"湘潭种莲历史悠久，所产莲子颗粒饱满、色白易煮，被称为“湘莲”，是传统出口农产品之一。夏季荷塘连片，秋后采莲、去壳、通心、晾晒，形成一整套手工技艺。当地流传“湘潭莲子不用淘”的说法，既是对品质的自信，也是几代种莲人积累的经验。今天的湘莲已经成为湘潭乡村特色产业的重要名片。",
			"/static/red/story-xianglian.png", "/static/red/ar-xianglian.png", 615},
	}
}

func seedFarmers() []seedFarmer {
	return []seedFarmer{
		{"farmer-xianglian", "杏花村湘莲合作社", "village-xingtang", "湘潭县白石镇杏花村",
			"由 32 户莲农组成的合作社，负责湘莲种植、通心加工与礼盒包装，为村内留守妇女提供稳定手工岗位。",
			"/static/red/farmer-lotus.png", "湘莲红色主题礼盒与通心白莲加工", 1200, 386},
		{"farmer-bamboo", "乌石竹编工坊", "village-wushi", "湘潭县乌石镇乌石村",
			"由老手艺人带队的竹编工坊，承担文创竹编包装与研学体验课程，订单按件计酬分给村民。",
			"/static/red/farmer-bamboo.png", "竹编包装与竹编手作课程", 800, 245},
		{"farmer-tea", "银田茶油作坊", "village-yinshan", "韶山市银田镇银田村",
			"传统木榨茶油作坊，联合周边农户收购油茶籽，采用统一包装与可追溯标识。",
			"/static/red/farmer-tea.png", "茶油红色文化伴手礼", 1000, 418},
		{"farmer-pepper", "泉湖剁辣椒小组", "village-quanhu", "湘乡市龙洞镇泉湖村",
			"由村内妇女组成的辣椒加工小组，负责剁辣椒、腊味的制作与灌装，配套统一的文创包装设计。",
			"/static/red/farmer-pepper.png", "剁辣椒与腊味伴手礼", 600, 302},
	}
}

func seedProducts() []seedProduct {
	return []seedProduct{
		{"product-bag-poem", "毛泽东诗词书法帆布袋", "红色经典IP",
			"/static/red/product-bag.png",
			"以《沁园春·长沙》词句为主视觉的加厚帆布袋，可日常通勤使用。",
			"选用 12 安加厚帆布，丝网印刷工艺，内衬加设小口袋。图案取自毛泽东同志诗词手迹的笔画结构，做当代化排布，兼顾纪念意义与日常实用性。由村内工坊承接后道缝制与包装，订单按件结算。",
			"《沁园春·长沙》写于一九二五年，词中“问苍茫大地，谁主沉浮”一句，是青年毛泽东面对国家命运的追问。",
			"farmer-bamboo", "乌石竹编工坊", "village-wushi",
			5800, 7800, 200, 126, []string{"帆布袋", "诗词", "日常可用"}, true, 1},
		{"product-badge-peng", "彭德怀“横刀立马”金属徽章", "红色经典IP",
			"/static/red/product-badge.png",
			"以“横刀立马”意象设计的金属烤漆徽章，可作胸针或包挂。",
			"锌合金压铸，双色烤漆工艺，背面为安全别针与挂环两用结构。造型提取自彭德怀同志纪念雕像的轮廓，做几何化处理，适合日常佩戴与收藏。",
			"毛泽东同志曾写下“谁敢横刀立马，唯我彭大将军”，这句评价让乌石村与彭德怀同志的军事形象紧密相连。",
			"", "", "village-wushi",
			2900, 3900, 300, 208, []string{"徽章", "元帅故里", "收藏"}, false, 2},
		{"product-sticker-guju", "韶山故居线稿贴纸套装", "红色经典IP",
			"/static/red/product-sticker.png",
			"六张一线条手绘的故居与纪念建筑贴纸，适合手账与笔记本装饰。",
			"采用防水哑光合成纸，模切边缘平整。图案由团队实地写生后绘制，包含故居、铜像广场、纪念馆、祠堂等六个场景，配套一张讲解卡片，扫码可读建筑背后的故事。",
			"贴纸上的每一座建筑，都对应着韶山革命纪念序列中的一个节点，从生活场景到纪念空间依次展开。",
			"", "", "village-shaoshan",
			1800, 2600, 500, 342, []string{"贴纸", "手绘", "手账"}, false, 3},
		{"product-blindbox", "红军长征主题盲盒手办", "红色经典IP",
			"/static/red/product-blindbox.png",
			"以长征经典场景为原型的小比例手办，共六款，含隐藏款。",
			"树脂材质，手工上色，高度约 8 厘米。六个常规款分别对应出发、行军、渡江、翻山、会师与宿营等场景，另设一款隐藏款。适合作为研学活动的激励奖品。",
			"长征是人类历史上的伟大远征，把手办做成场景序列，是希望用年轻人熟悉的方式讲清这段历程。",
			"", "", "village-yinshan",
			6900, 8900, 150, 87, []string{"盲盒", "手办", "研学奖品"}, false, 4},
		{"product-scarf-shrimp", "齐白石虾图真丝丝巾", "齐白石艺术",
			"/static/red/product-scarf.png",
			"取齐白石《虾》作品意趣设计的真丝方巾，献给喜欢传统艺术的你。",
			"桑蚕丝斜纹绸，手工卷边，尺寸 90×90 厘米。图案在保留原作笔墨节奏的基础上重新构图，适配日常搭配与礼品场景。",
			"齐白石画虾数十年，从形似到神似，最终以极简笔墨写出一池清趣，是“衰年变法”的代表题材。",
			"", "", "village-xingtang",
			29800, 36800, 60, 34, []string{"真丝", "齐白石", "礼品"}, false, 5},
		{"product-seal", "齐白石篆刻印章文创礼盒", "齐白石艺术",
			"/static/red/product-seal.png",
			"含一方定制印章与印泥的礼盒，可刻姓名或字号。",
			"礼盒含青田石印材一方、仿古印泥一盒、红木印盒一只、印谱册一册。印章由本地手艺人手工篆刻，下单后按需求刻制，适合作为升学或入职礼物。",
			"齐白石早年以雕花木匠为生，后精于篆刻，其“单刀直入”的刀法自成一格。",
			"", "", "village-xingtang",
			39800, 49800, 40, 21, []string{"篆刻", "定制", "礼盒"}, true, 6},
		{"product-postcard", "星斗塘老屋水彩明信片套装", "齐白石艺术",
			"/static/red/product-postcard.png",
			"八张水彩明信片，描绘星斗塘周边的老屋、荷塘与田埂。",
			"由团队成员实地写生后手绘上色，选用 300 克特种纸，附赠可寄递信封。每张明信片背面印有二维码，扫描后可看到对应的语音讲解与老照片。",
			"星斗塘是齐白石出生的地方，荷塘与田埂构成了他最早的美术启蒙环境。",
			"", "", "village-xingtang",
			2800, 3600, 400, 256, []string{"明信片", "扫码故事", "AR"}, false, 7},
		{"product-book-baishi", "“从木匠到大师”成长历程绘本", "齐白石艺术",
			"/static/red/product-book.png",
			"面向青少年的图文绘本，讲述齐白石从木匠到画坛宗师的成长历程。",
			"精装 48 页，全彩印刷，由团队插画组绘制。内容分为学木匠、遇师友、走南北、衰年变法四段，配套一份研学任务单，适合中小学使用。",
			"“不教一日闲过”是齐白石一生勤奋的写照，这本书希望把这种精神讲给孩子听。",
			"", "", "village-xingtang",
			5800, 6800, 200, 143, []string{"绘本", "青少年", "研学"}, false, 8},
		{"product-lotus-gift", "湘莲红色主题包装礼盒", "湘村风物助农",
			"/static/red/product-lotus.png",
			"通心白莲礼盒，配红色文化主题包装，产地直发。",
			"净含量 500 克，选用湘潭本地当年新莲，人工通心、自然晾晒。包装由团队设计，提取荷花与红旗元素，盒内附合作社溯源信息卡与食用建议。",
			"湘潭种莲历史悠久，湘莲颗粒饱满、易煮易糯，是几代湘潭人的餐桌记忆。",
			"farmer-xianglian", "杏花村湘莲合作社", "village-xingtang",
			8800, 10800, 300, 178, []string{"助农", "湘莲", "产地直发"}, false, 9},
		{"product-camellia-oil", "湘潭茶油红色文化伴手礼", "湘村风物助农",
			"/static/red/product-camellia.png",
			"传统木榨茶油双瓶装，配文创外盒，适合探亲与商务送礼。",
			"每瓶 500 毫升，低温压榨与物理精滤工艺，酸值与过氧化值符合国家标准。外盒采用可回收纸材，印有银田灌区与油茶林插画。",
			"茶油是湖南传统食用油料，银田一带的油茶林多为农户分散种植，统一包装后附加值明显提升。",
			"farmer-tea", "银田茶油作坊", "village-yinshan",
			19800, 23800, 150, 96, []string{"助农", "茶油", "伴手礼"}, false, 10},
		{"product-bamboo-craft", "乌石峰“三十六坊”竹编手作", "湘村风物助农",
			"/static/red/product-bamboo.png",
			"老手艺人手工编制的竹编收纳篮，可用于文创包装与家居收纳。",
			"选用三年生毛竹，经破篾、刮青、编织、蒸煮定型等工序完成。尺寸约 22×16×12 厘米，每件略有差异。订单由工坊按件分派给村民，收益直接进入农户家庭。",
			"竹编是湘中乡村的传统手艺，一只篮子往往要经十几道工序，全部依靠双手完成。",
			"farmer-bamboo", "乌石竹编工坊", "village-wushi",
			12800, 15800, 120, 65, []string{"助农", "非遗手作", "竹编"}, true, 11},
		{"product-chili-sauce", "泉湖剁辣椒与腊味伴手礼", "湘村风物助农",
			"/static/red/product-chili.png",
			"村内妇女小组制作的剁辣椒与腊味组合装，湘味十足。",
			"剁辣椒 2 瓶（每瓶 300 克）、腊肉 1 份（500 克），真空包装，冷链或常温按季节发货。外盒印有村落故事卡，扫码可听制作人的语音介绍。",
			"剁辣椒是湘潭人餐桌上的常客，村内妇女小组用手工剁制的方式保留了传统口味。",
			"farmer-pepper", "泉湖剁辣椒小组", "village-quanhu",
			13800, 16800, 100, 58, []string{"助农", "湘味", "组合装"}, false, 12},
		{"product-digital-card", "红色村落数字纪念卡", "数字文创",
			"/static/red/product-digital.png",
			"含六处红色村落数字纪念卡，附带专属编号与收藏证书。",
			"数字纪念卡以六处红色村落的标志性建筑为原型设计，购买后可在线收藏、查看编号与创作者手记，并可在小程序内展示。配套纸质收藏证书一张，可随商品寄出。",
			"把数字纪念卡绑定到真实村落，是希望每一次收藏都能对应到一片土地上的记忆。",
			"", "", "village-shaoshan",
			3900, 4900, 999, 412, []string{"数字文创", "纪念卡", "收藏"}, false, 13},
	}
}

func seedTasks() []seedTask {
	return []seedTask{
		{"task-maoguju", "village-shaoshan", "spot-maoguju", "在故居寻找三处生活细节", "现场探访",
			"走进毛泽东同志故居，找出堂屋、灶房与卧室的位置关系，记录一处让你印象最深的陈设。", "/static/red/task-maoguju.png", 20, 1},
		{"task-nongyun", "village-shaoshan", "spot-maozhendang", "还原一段农民运动史", "史料研读",
			"在纪念馆中找到 1925 年韶山农民运动的展陈内容，用自己的话写下 100 字以内的故事梗概。", "/static/red/task-nongyun.png", 30, 2},
		{"task-hengdao", "village-wushi", "spot-pengjinianguan", "写下你理解的“横刀立马”", "主题写作",
			"参观彭德怀纪念馆后，用一句话写下你对“横刀立马”这种精神的理解，并说明理由。", "/static/red/task-hengdao.png", 30, 3},
		{"task-baishi", "village-xingtang", "spot-xingdoutang", "临摹一只虾", "写生实践",
			"在星斗塘观察荷塘景致，参考齐白石的用笔方式，用十分钟完成一幅速写或临摹。", "/static/red/task-baishi.png", 40, 4},
		{"task-guanqu", "village-yinshan", "spot-guanqu", "算一算灌区的水", "数据调查",
			"了解韶山灌区的干渠长度与灌溉面积，用简单图表记录这组数字，并思考它对当地农业的意义。", "/static/red/task-guanqu.png", 30, 5},
		{"task-chengeng", "village-quanhu", "spot-chengengguju", "将军生平时间轴", "小组协作",
			"小组合作梳理陈赓同志的重要经历，做成一条时间轴，在故居前做 2 分钟讲解。", "/static/red/task-chengeng.png", 30, 6},
	}
}

func (s *MySQLStore) seedRedTourism() error {
	ctx := context.Background()

	for _, item := range seedVillages() {
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO villages
			(id,name,region,summary,history,cover_url,lat,lng,honors,themes,spot_count,story_count,product_count,sort)
			VALUES (?,?,?,?,?,?,?,?,?,?,0,0,0,?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), region=VALUES(region), summary=VALUES(summary),
			history=VALUES(history), cover_url=VALUES(cover_url), lat=VALUES(lat), lng=VALUES(lng),
			honors=VALUES(honors), themes=VALUES(themes), sort=VALUES(sort)`,
			item.id, item.name, item.region, item.summary, item.history, item.cover, item.lat, item.lng,
			jsonList(item.honors), jsonList(item.themes), item.sort); err != nil {
			return err
		}
	}
	for _, item := range seedSpots() {
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO spots
			(id,village_id,name,category,summary,image_url,audio_url,lat,lng,sort)
			VALUES (?,?,?,?,?,?,'',?,?,?)
			ON DUPLICATE KEY UPDATE village_id=VALUES(village_id), name=VALUES(name), category=VALUES(category),
			summary=VALUES(summary), image_url=VALUES(image_url), lat=VALUES(lat), lng=VALUES(lng), sort=VALUES(sort)`,
			item.id, item.villageID, item.name, item.category, item.summary, item.image, item.lat, item.lng, item.sort); err != nil {
			return err
		}
	}
	for _, item := range seedRoutes() {
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO routes
			(id,name,theme,summary,cover_url,days,distance_km,spot_ids,spot_names,sort)
			VALUES (?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), theme=VALUES(theme), summary=VALUES(summary),
			cover_url=VALUES(cover_url), days=VALUES(days), distance_km=VALUES(distance_km),
			spot_ids=VALUES(spot_ids), spot_names=VALUES(spot_names), sort=VALUES(sort)`,
			item.id, item.name, item.theme, item.summary, item.cover, item.days, item.distance,
			jsonList(item.spotIDs), jsonList(item.spotNames), item.sort); err != nil {
			return err
		}
	}
	for _, item := range seedStories() {
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO stories
			(id,village_id,title,category,era,summary,content,cover_url,audio_url,ar_image,views)
			VALUES (?,?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE village_id=VALUES(village_id), title=VALUES(title), category=VALUES(category),
			era=VALUES(era), summary=VALUES(summary), content=VALUES(content), cover_url=VALUES(cover_url),
			ar_image=VALUES(ar_image)`,
			item.id, item.villageID, item.title, item.category, item.era, item.summary, item.content,
			item.cover, "", item.arImage, item.views); err != nil {
			return err
		}
	}
	for _, item := range seedFarmers() {
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO farmers
			(id,name,village_id,village_name,intro,avatar_url,project,support_target,support_count)
			VALUES (?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), village_id=VALUES(village_id), village_name=VALUES(village_name),
			intro=VALUES(intro), avatar_url=VALUES(avatar_url), project=VALUES(project),
			support_target=VALUES(support_target)`,
			item.id, item.name, item.villageID, item.villageName, item.intro, item.avatar, item.project,
			item.target, item.count); err != nil {
			return err
		}
	}
	for _, item := range seedProducts() {
		custom := 0
		if item.customizable {
			custom = 1
		}
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO products
			(id,name,category,price_cents,points_price,reward_points,reward_exp,original_cents,stock,sales,cover_url,images,summary,detail,story_card,
			tags,customizable,farmer_id,farmer_name,village_id,status,sort)
			VALUES (?,?,?,?,?,?,?, ?,?,?,?,'[]',?,?,?,?,?,?,?,?,'上架',?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), category=VALUES(category), price_cents=VALUES(price_cents),
			points_price=VALUES(points_price), reward_points=VALUES(reward_points), reward_exp=VALUES(reward_exp),
			original_cents=VALUES(original_cents), stock=VALUES(stock), cover_url=VALUES(cover_url),
			summary=VALUES(summary), detail=VALUES(detail), story_card=VALUES(story_card), tags=VALUES(tags),
			customizable=VALUES(customizable), farmer_id=VALUES(farmer_id), farmer_name=VALUES(farmer_name),
			village_id=VALUES(village_id), sort=VALUES(sort)`,
			item.id, item.name, item.category, item.price, defaultPointsPrice(item.price), defaultRewardPoints(item.price), defaultRewardExp(item.price),
			item.original, item.stock, item.sales, item.cover, item.summary, item.detail, item.storyCard, jsonList(item.tags), custom,
			item.farmerID, item.farmerName, item.villageID, item.sort); err != nil {
			return err
		}
	}
	for _, item := range seedTasks() {
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO tasks
			(id,village_id,spot_id,title,type,summary,points,reward_exp,cover_url,sort)
			VALUES (?,?,?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE village_id=VALUES(village_id), spot_id=VALUES(spot_id), title=VALUES(title),
			type=VALUES(type), summary=VALUES(summary), points=VALUES(points), reward_exp=VALUES(reward_exp), cover_url=VALUES(cover_url), sort=VALUES(sort)`,
			item.id, item.villageID, item.spotID, item.title, item.kind, item.summary, item.points, item.points, item.cover, item.sort); err != nil {
			return err
		}
	}
	banners := []Banner{
		{"banner-village", "红映湘村", "湘潭红色村落数字文旅与文创助农", "/static/red/banner-village.png", "village", "village-shaoshan", 1, true},
		{"banner-market", "红文创商城", "文创好物与助农特产，每一单都分给村里", "/static/red/banner-market.png", "tab", "mall", 2, true},
		{"banner-story", "红迹数字故事", "扫码读故事，听湘潭红色村落的记忆", "/static/red/banner-story.png", "tab", "story", 3, true},
	}
	for _, item := range banners {
		active := 0
		if item.Active {
			active = 1
		}
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO banners
			(id,title,subtitle,image_url,link_type,link_id,sort,active) VALUES (?,?,?,?,?,?,?,?)
			ON DUPLICATE KEY UPDATE title=VALUES(title), subtitle=VALUES(subtitle), image_url=VALUES(image_url),
			link_type=VALUES(link_type), link_id=VALUES(link_id), sort=VALUES(sort), active=VALUES(active)`,
			item.ID, item.Title, item.Subtitle, item.ImageURL, item.LinkType, item.LinkID, item.Sort, active); err != nil {
			return err
		}
	}
	return s.syncVillageCounts(ctx)
}
