package constants

import "qrpay/utils"

type VietQRStatus int

const (
	VietQRNotSupported      VietQRStatus = -1
	VietQRReceiveOnly       VietQRStatus = 0
	VietQRTransferSupported VietQRStatus = 1
)

type Bank struct {
	Key             BankKey      `json:"Key"`
	Code            BankCode     `json:"code"`
	Name            string       `json:"name"`
	ShortName       string       `json:"shortName"`
	BIN             string       `json:"bin"`
	VietQRStatus    VietQRStatus `json:"vietQRStatus"`
	LookupSupported *int         `json:"lookupSupported,omitempty"`
	SWIFTCode       *string      `json:"swiftCode,omitempty"`
	Keywords        *string      `json:"Keywords,omitempty"`
	Deprecated      *bool        `json:"deprecated,omitempty"`
}

// BanksMap maps bank Keys to Bank structs.
var BanksMap = map[BankKey]Bank{
	BankKey_ABBANK: {
		Key:             "ABBANK",
		Code:            "ABBANK",
		Name:            "Ngân hàng TMCP An Bình",
		BIN:             "970425",
		ShortName:       "AB Bank",
		VietQRStatus:    VietQRTransferSupported,
		LookupSupported: utils.ToPointer(1),
		SWIFTCode:       utils.ToPointer("ABBKVNVX"),
		Keywords:        utils.ToPointer("anbinh"),
	},
	BankKey_ACB: {
		Key:             "ACB",
		Code:            "ACB",
		Name:            "Ngân hàng TMCP Á Châu",
		BIN:             "970416",
		ShortName:       "ACB",
		VietQRStatus:    VietQRTransferSupported,
		LookupSupported: utils.ToPointer(1),
		SWIFTCode:       utils.ToPointer("ASCBVNVX"),
		Keywords:        utils.ToPointer("achau"),
	},
BankKey_AGRIBANK: {
	Key: BankKey_AGRIBANK,
	Code: BankCode_AGRIBANK,
	Name: 'Ngân hàng Nông nghiệp và Phát triển Nông thôn Việt Nam',
	BIN: '970405',
	ShortName: 'Agribank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VBAAVNVX',
	Keywords: 'nongnghiep, nongthon, agribank, agri'
},
BankKey_BAC_A_BANK: {
	Key: BankKey_BAC_A_BANK,
	Code: BankCode_BAC_A_BANK,
	Name: 'Ngân hàng TMCP Bắc Á',
	BIN: '970409',
	ShortName: 'BacA Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'NASCVNVX',
	Keywords: 'baca, NASB'
},
BankKey_BAOVIET_BANK: {
	Key: BankKey_BAOVIET_BANK,
	Code: BankCode_BAOVIET_BANK,
	Name: 'Ngân hàng TMCP Bảo Việt',
	BIN: '970438',
	ShortName: 'BaoViet Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'BVBVVNVX',
	Keywords: 'baoviet, BVB'
},
BankKey_BANVIET: {
	Key: BankKey_BANVIET,
	Code: BankCode_BANVIET,
	Name: 'Ngân hàng TMCP Bản Việt',
	BIN: '970454',
	ShortName: 'BVBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VCBCVNVX',
	Keywords: 'banviet, vietcapitalbank'
},
BankKey_BIDC: {
	Key: BankKey_BIDC,
	Code: BankCode_BIDC,
	Name: 'Ngân hàng TMCP Đầu tư và Phát triển Campuchia',
	BIN: '',
	ShortName: 'BIDC',
	VietQRStatus: VietQRStatus.NOT_SUPPORTED
},
BankKey_BIDV: {
	Key: BankKey_BIDV,
	Code: BankCode_BIDV,
	Name: 'Ngân hàng TMCP Đầu tư và Phát triển Việt Nam',
	BIN: '970418',
	ShortName: 'BIDV',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'BIDVVNVX'
},
BankKey_CAKE: {
	Key: BankKey_CAKE,
	Code: BankCode_CAKE,
	Name: 'Ngân hàng số CAKE by VPBank - Ngân hàng TMCP Việt Nam Thịnh Vượng',
	BIN: '546034',
	ShortName: 'CAKE by VPBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: null
},
BankKey_CBBANK: {
	Key: BankKey_CBBANK,
	Code: BankCode_CBBANK,
	Name: 'Ngân hàng Thương mại TNHH MTV Xây dựng Việt Nam',
	BIN: '970444',
	ShortName: 'CB Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'GTBAVNVX',
	Keywords: 'xaydungvn, xaydung'
},
BankKey_CIMB: {
	Key: BankKey_CIMB,
	Code: BankCode_CIMB,
	Name: 'Ngân hàng TNHH MTV CIMB Việt Nam',
	BIN: '422589',
	ShortName: 'CIMB Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'CIBBVNVN',
	Keywords: 'cimbvn'
},
BankKey_COOP_BANK: {
	Key: BankKey_COOP_BANK,
	Code: BankCode_COOP_BANK,
	Name: 'Ngân hàng Hợp tác xã Việt Nam',
	BIN: '970446',
	ShortName: 'Co-op Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: null,
	Keywords: 'hoptacxa, coop'
},
BankKey_DBS_BANK: {
	Key: BankKey_DBS_BANK,
	Code: BankCode_DBS_BANK,
	Name: 'NH TNHH MTV Phát triển Singapore - Chi nhánh TP. Hồ Chí Minh',
	BIN: '796500',
	ShortName: 'DBS Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 0,
	swiftCode: 'DBSSVNVX',
	Keywords: 'dbshcm'
},
BankKey_DONG_A_BANK: {
	Key: BankKey_DONG_A_BANK,
	Code: BankCode_DONG_A_BANK,
	Name: 'Ngân hàng TMCP Đông Á',
	BIN: '970406',
	ShortName: 'DongA Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'EACBVNVX',
	Keywords: 'donga, DAB'
},
BankKey_EXIMBANK: {
	Key: BankKey_EXIMBANK,
	Code: BankCode_EXIMBANK,
	Name: 'Ngân hàng TMCP Xuất Nhập khẩu Việt Nam',
	BIN: '970431',
	ShortName: 'Eximbank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'EBVIVNVX'
},
BankKey_GPBANK: {
	Key: BankKey_GPBANK,
	Code: BankCode_GPBANK,
	Name: 'Ngân hàng Thương mại TNHH MTV Dầu Khí Toàn Cầu',
	BIN: '970408',
	ShortName: 'GPBank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'GBNKVNVX',
	Keywords: 'daukhi'
},
BankKey_HDBANK: {
	Key: BankKey_HDBANK,
	Code: BankCode_HDBANK,
	Name: 'Ngân hàng TMCP Phát triển TP. Hồ Chí Minh',
	BIN: '970437',
	ShortName: 'HDBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'HDBCVNVX'
},
BankKey_HONGLEONG_BANK: {
	Key: BankKey_HONGLEONG_BANK,
	Code: BankCode_HONGLEONG_BANK,
	Name: 'Ngân hàng TNHH MTV Hong Leong Việt Nam',
	BIN: '970442',
	ShortName: 'HongLeong Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'HLBBVNVX',
	Keywords: 'HLBVN'
},
BankKey_HSBC: {
	Key: BankKey_HSBC,
	Code: BankCode_HSBC,
	Name: 'Ngân hàng TNHH MTV HSBC (Việt Nam)',
	BIN: '458761',
	ShortName: 'HSBC Vietnam',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'HSBCVNVX'
},
BankKey_IBK_HCM: {
	Key: BankKey_IBK_HCM,
	Code: BankCode_IBK_HCM,
	Name: 'Ngân hàng Công nghiệp Hàn Quốc - Chi nhánh TP. Hồ Chí Minh',
	BIN: '970456',
	ShortName: 'IBK HCM',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 0,
	swiftCode: null,
	Keywords: 'congnghiep'
},
BankKey_IBK_HN: {
	Key: BankKey_IBK_HN,
	Code: BankCode_IBK_HN,
	Name: 'Ngân hàng Công nghiệp Hàn Quốc - Chi nhánh Hà Nội',
	BIN: '970455',
	ShortName: 'IBK HN',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 0,
	swiftCode: null,
	Keywords: 'congnghiep'
},
BankKey_INDOVINA_BANK: {
	Key: BankKey_INDOVINA_BANK,
	Code: BankCode_INDOVINA_BANK,
	Name: 'Ngân hàng TNHH Indovina',
	BIN: '970434',
	ShortName: 'Indovina Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: null
},
BankKey_KASIKORN_BANK: {
	Key: BankKey_KASIKORN_BANK,
	Code: BankCode_KASIKORN_BANK,
	Name: 'Ngân hàng Đại chúng TNHH KASIKORNBANK - CN TP. Hồ Chí Minh',
	BIN: '668888',
	ShortName: 'Kasikornbank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'KASIVNVX'
},
BankKey_KIENLONG_BANK: {
	Key: BankKey_KIENLONG_BANK,
	Code: BankCode_KIENLONG_BANK,
	Name: 'Ngân hàng TMCP Kiên Long',
	BIN: '970452',
	ShortName: 'KienlongBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'KLBKVNVX'
},
BankKey_KOOKMIN_BANK_HCM: {
	Key: BankKey_KOOKMIN_BANK_HCM,
	Code: BankCode_KOOKMIN_BANK_HCM,
	Name: 'Ngân hàng Kookmin - Chi nhánh TP. Hồ Chí Minh',
	BIN: '970463',
	ShortName: 'Kookmin Bank HCM',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 0,
	swiftCode: null
},
BankKey_KOOKMIN_BANK_HN: {
	Key: BankKey_KOOKMIN_BANK_HN,
	Code: BankCode_KOOKMIN_BANK_HN,
	Name: 'Ngân hàng Kookmin - Chi nhánh Hà Nội',
	BIN: '970462',
	ShortName: 'Kookmin Bank HN',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 0,
	swiftCode: null
},
BankKey_LIENVIETPOST_BANK: {
	Key: BankKey_LIENVIETPOST_BANK,
	Code: BankCode_LPBANK,
	Name: 'Ngân hàng TMCP Bưu Điện Liên Việt',
	BIN: '970449',
	ShortName: 'LienVietPostBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'LVBKVNVX',
	Keywords: 'lienvietbank',
	deprecated: true
},
BankKey_LPBANK: {
	Key: BankKey_LPBANK,
	Code: BankCode_LPBANK,
	Name: 'Ngân hàng TMCP Lộc Phát Việt Nam',
	BIN: '970449',
	ShortName: 'LPBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'LVBKVNVX',
	Keywords: 'lienvietbank, loc phat',
},
BankKey_LIOBANK: {
	Key: BankKey_LIOBANK,
	Code: BankCode_LIOBANK,
	Name: 'Ngân hàng số Liobank - Ngân hàng TMCP Phương Đông',
	BIN: '963369',
	ShortName: 'Liobank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: null
},
BankKey_MBBANK: {
	Key: BankKey_MBBANK,
	Code: BankCode_MBBANK,
	Name: 'Ngân hàng TMCP Quân đội',
	BIN: '970422',
	ShortName: 'MB Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'MSCBVNVX'
},
BankKey_MBV: {
	Key: BankKey_MBV,
	Code: BankCode_MBV,
	Name: 'Ngân hàng TNHH MTV Việt Nam Hiện Đại',
	BIN: '970414',
	ShortName: 'MBV',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'OCBKUS3M',
	Keywords: 'daiduong, mbv',
},
BankKey_MSB: {
	Key: BankKey_MSB,
	Code: BankCode_MSB,
	Name: 'Ngân hàng TMCP Hàng Hải',
	BIN: '970426',
	ShortName: 'MSB',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'MCOBVNVX',
	Keywords: 'hanghai'
},
BankKey_NAM_A_BANK: {
	Key: BankKey_NAM_A_BANK,
	Code: BankCode_NAM_A_BANK,
	Name: 'Ngân hàng TMCP Nam Á',
	BIN: '970428',
	ShortName: 'Nam A Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'NAMAVNVX',
	Keywords: 'namabank'
},
BankKey_NCB: {
	Key: BankKey_NCB,
	Code: BankCode_NCB,
	Name: 'Ngân hàng TMCP Quốc Dân',
	BIN: '970419',
	ShortName: 'NCB Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'NVBAVNVX',
	Keywords: 'quocdan'
},
BankKey_NONGHYUP_BANK_HN: {
	Key: BankKey_NONGHYUP_BANK_HN,
	Code: BankCode_NONGHYUP_BANK_HN,
	Name: 'Ngân hàng Nonghyup - Chi nhánh Hà Nội',
	BIN: '801011',
	ShortName: 'Nonghyup Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 0,
	swiftCode: null
},
BankKey_OCB: {
	Key: BankKey_OCB,
	Code: BankCode_OCB,
	Name: 'Ngân hàng TMCP Phương Đông',
	BIN: '970448',
	ShortName: 'OCB Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'ORCOVNVX',
	Keywords: 'phuongdong'
},
BankKey_OCEANBANK: {
	Key: BankKey_OCEANBANK,
	Code: BankCode_OCEANBANK,
	Name: 'Ngân hàng Thương mại TNHH MTV Đại Dương',
	BIN: '970414',
	ShortName: 'Ocean Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'OCBKUS3M',
	Keywords: 'daiduong',
	deprecated: true
},
BankKey_PGBANK: {
	Key: BankKey_PGBANK,
	Code: BankCode_PGBANK,
	Name: 'Ngân hàng TMCP Xăng dầu Petrolimex',
	BIN: '970430',
	ShortName: 'PG Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'PGBLVNVX'
},
BankKey_PUBLIC_BANK: {
	Key: BankKey_PUBLIC_BANK,
	Code: BankCode_PUBLIC_BANK,
	Name: 'Ngân hàng TNHH MTV Public Việt Nam',
	BIN: '970439',
	ShortName: 'Public Bank Vietnam',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'VIDPVNVX',
	Keywords: 'publicvn'
},
BankKey_PVCOM_BANK: {
	Key: BankKey_PVCOM_BANK,
	Code: BankCode_PVCOM_BANK,
	Name: 'Ngân hàng TMCP Đại Chúng Việt Nam',
	BIN: '970412',
	ShortName: 'PVcomBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'WBVNVNVX',
	Keywords: 'daichung'
},
BankKey_SACOMBANK: {
	Key: BankKey_SACOMBANK,
	Code: BankCode_SACOMBANK,
	Name: 'Ngân hàng TMCP Sài Gòn Thương Tín',
	BIN: '970403',
	ShortName: 'Sacombank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'SGTTVNVX'
},
BankKey_SAIGONBANK: {
	Key: BankKey_SAIGONBANK,
	Code: BankCode_SAIGONBANK,
	Name: 'Ngân hàng TMCP Sài Gòn Công Thương',
	BIN: '970400',
	ShortName: 'SaigonBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'SBITVNVX',
	Keywords: 'saigoncongthuong, saigonbank'
},
BankKey_SCB: {
	Key: BankKey_SCB,
	Code: BankCode_SCB,
	Name: 'Ngân hàng TMCP Sài Gòn',
	BIN: '970429',
	ShortName: 'SCB',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'SACLVNVX',
	Keywords: 'saigon'
},
BankKey_SEA_BANK: {
	Key: BankKey_SEA_BANK,
	Code: BankCode_SEA_BANK,
	Name: 'Ngân hàng TMCP Đông Nam Á',
	BIN: '970440',
	ShortName: 'SeABank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'SEAVVNVX'
},
BankKey_SHB: {
	Key: BankKey_SHB,
	Code: BankCode_SHB,
	Name: 'Ngân hàng TMCP Sài Gòn - Hà Nội',
	BIN: '970443',
	ShortName: 'SHB',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'SHBAVNVX',
	Keywords: 'saigonhanoi, sghn'
},
BankKey_SHINHAN_BANK: {
	Key: BankKey_SHINHAN_BANK,
	Code: BankCode_SHINHAN_BANK,
	Name: 'Ngân hàng TNHH MTV Shinhan Việt Nam',
	BIN: '970424',
	ShortName: 'Shinhan Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'SHBKVNVX'
},
BankKey_STANDARD_CHARTERED_BANK: {
	Key: BankKey_STANDARD_CHARTERED_BANK,
	Code: BankCode_STANDARD_CHARTERED_BANK,
	Name: 'Ngân hàng TNHH MTV Standard Chartered Bank Việt Nam',
	BIN: '970410',
	ShortName: 'Standard Chartered Vietnam',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: 'SCBLVNVX'
},
BankKey_TECHCOMBANK: {
	Key: BankKey_TECHCOMBANK,
	Code: BankCode_TECHCOMBANK,
	Name: 'Ngân hàng TMCP Kỹ thương Việt Nam',
	BIN: '970407',
	ShortName: 'Techcombank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VTCBVNVX'
},
BankKey_TIMO: {
	Key: BankKey_TIMO,
	Code: BankCode_TIMO,
	Name: 'Ngân hàng số Timo by Bản Việt Bank',
	BIN: '963388',
	ShortName: 'Timo',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 0,
	swiftCode: null,
	Keywords: 'banviet'
},
BankKey_TPBANK: {
	Key: BankKey_TPBANK,
	Code: BankCode_TPBANK,
	Name: 'Ngân hàng TMCP Tiên Phong',
	BIN: '970423',
	ShortName: 'TPBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'TPBVVNVX',
	Keywords: 'tienphong'
},
BankKey_UBANK: {
	Key: BankKey_UBANK,
	Code: BankCode_UBANK,
	Name: 'Ngân hàng số Ubank by VPBank - Ngân hàng TMCP Việt Nam Thịnh Vượng',
	BIN: '546035',
	ShortName: 'Ubank by VPBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: null
},
BankKey_UNITED_OVERSEAS_BANK: {
	Key: BankKey_UNITED_OVERSEAS_BANK,
	Code: BankCode_UNITED_OVERSEAS_BANK,
	Name: 'Ngân hàng United Overseas Bank Việt Nam',
	BIN: '970458',
	ShortName: 'United Overseas Bank Vietnam',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: null
},
BankKey_VIB: {
	Key: BankKey_VIB,
	Code: BankCode_VIB,
	Name: 'Ngân hàng TMCP Quốc tế Việt Nam',
	BIN: '970441',
	ShortName: 'VIB',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VNIBVNVX',
	Keywords: 'quocte'
},
BankKey_VIET_A_BANK: {
	Key: BankKey_VIET_A_BANK,
	Code: BankCode_VIET_A_BANK,
	Name: 'Ngân hàng TMCP Việt Á',
	BIN: '970427',
	ShortName: 'VietABank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VNACVNVX'
},
BankKey_VIET_BANK: {
	Key: BankKey_VIET_BANK,
	Code: BankCode_VIET_BANK,
	Name: 'Ngân hàng TMCP Việt Nam Thương Tín',
	BIN: '970433',
	ShortName: 'VietBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VNTTVNVX',
	Keywords: 'vietnamthuongtin, vnthuongtin'
},
BankKey_VIETCOMBANK: {
	Key: BankKey_VIETCOMBANK,
	Code: BankCode_VIETCOMBANK,
	Name: 'Ngân hàng TMCP Ngoại Thương Việt Nam',
	BIN: '970436',
	ShortName: 'Vietcombank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'BFTVVNVX'
},
BankKey_VIETINBANK: {
	Key: BankKey_VIETINBANK,
	Code: BankCode_VIETINBANK,
	Name: 'Ngân hàng TMCP Công thương Việt Nam',
	BIN: '970415',
	ShortName: 'VietinBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'ICBVVNVX',
	Keywords: 'viettin' // Some users may use this Keyword
},
BankKey_VIKKI: {
	Key: BankKey_VIKKI,
	Code: BankCode_VIKKI,
	Name: 'Ngân hàng TNHH MTV Số Vikki',
	BIN: '970406',
	ShortName: 'Vikki Bank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'EACBVNVX',
	Keywords: 'vikki, dongabank, dong a',
},
BankKey_VPBANK: {
	Key: BankKey_VPBANK,
	Code: BankCode_VPBANK,
	Name: 'Ngân hàng TMCP Việt Nam Thịnh Vượng',
	BIN: '970432',
	ShortName: 'VPBank',
	VietQRStatus: VietQRStatus.TRANSFER_SUPPORTED,
	lookupSupported: 1,
	swiftCode: 'VPBKVNVX',
	Keywords: 'vnthinhvuong'
},
BankKey_VRB: {
	Key: BankKey_VRB,
	Code: BankCode_VRB,
	Name: 'Ngân hàng Liên doanh Việt - Nga',
	BIN: '970421',
	ShortName: 'VietNgaBank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: null,
	Keywords: 'vietnam-russia, vrbank'
},
BankKey_WOORI_BANK: {
	Key: BankKey_WOORI_BANK,
	Code: BankCode_WOORI_BANK,
	Name: 'Ngân hàng TNHH MTV Woori Việt Nam',
	BIN: '970457',
	ShortName: 'Woori Bank',
	VietQRStatus: VietQRStatus.RECEIVE_ONLY,
	lookupSupported: 1,
	swiftCode: null
}
}

// Banks is a slice of all banks.
var Banks = func() []Bank {
	banks := make([]Bank, 0, len(BanksMap))
	for _, bank := range BanksMap {
		banks = append(banks, bank)
	}
	return banks
}()
