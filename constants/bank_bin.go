package constants

// BankBIN is the 6-digit Napas bank identifier carried in VietQR field 38.01.00.
type BankBIN string

func (b BankBIN) String() string { return string(b) }

// Common BINs. Use BinForKey / BinForCode for runtime lookup; use the
// constants below for compile-time references.
const (
	BIN_ABBANK      BankBIN = "970425" // An Binh Bank
	BIN_ACB         BankBIN = "970416" // ACB
	BIN_AGRIBANK    BankBIN = "970405" // Agribank
	BIN_BACABANK    BankBIN = "970409" // Bac A Bank
	BIN_BAOVIET     BankBIN = "970438" // BaoViet Bank
	BIN_BIDV        BankBIN = "970418" // BIDV
	BIN_BVBANK      BankBIN = "970454" // BVBank / BanViet (MoMo + ZaloPay tunnel)
	BIN_CAKE        BankBIN = "546034" // Cake by VPBank
	BIN_CBBANK      BankBIN = "970444" // CB Bank
	BIN_CIMB        BankBIN = "422589" // CIMB
	BIN_COOPBANK    BankBIN = "970446" // Co-op Bank
	BIN_DBSBANK     BankBIN = "796500" // DBS Bank
	BIN_DONGABANK   BankBIN = "970406" // DongA Bank (deprecated)
	BIN_EXIMBANK    BankBIN = "970431" // Eximbank
	BIN_GPBANK      BankBIN = "970408" // GPBank
	BIN_HDBANK      BankBIN = "970437" // HDBank
	BIN_HONGLEONG   BankBIN = "970442" // Hong Leong Bank
	BIN_HSBC        BankBIN = "458761" // HSBC
	BIN_IBKHCM      BankBIN = "970456" // IBK HCM
	BIN_IBKHN       BankBIN = "970455" // IBK HN
	BIN_INDOVINA    BankBIN = "970434" // Indovina Bank
	BIN_KASIKORN    BankBIN = "668888" // Kasikorn Bank
	BIN_KIENLONG    BankBIN = "970452" // KienLong Bank
	BIN_KOOKMINHCM  BankBIN = "970463" // Kookmin Bank HCM
	BIN_KOOKMINHN   BankBIN = "970462" // Kookmin Bank HN
	BIN_LPBANK      BankBIN = "970449" // LPBank
	BIN_MBBANK      BankBIN = "970422" // MB Bank
	BIN_MSB         BankBIN = "970426" // MSB
	BIN_NAMABANK    BankBIN = "970428" // Nam A Bank
	BIN_NCB         BankBIN = "970419" // NCB
	BIN_OCB         BankBIN = "970448" // OCB
	BIN_PGBANK      BankBIN = "970430" // PGBank
	BIN_PUBLICBANK  BankBIN = "970439" // Public Bank
	BIN_PVCOMBANK   BankBIN = "970412" // PVcomBank
	BIN_SACOMBANK   BankBIN = "970403" // Sacombank
	BIN_SAIGONBANK  BankBIN = "970400" // Saigonbank
	BIN_SCB         BankBIN = "970429" // SCB
	BIN_SEABANK     BankBIN = "970440" // SeABank
	BIN_SHB         BankBIN = "970443" // SHB
	BIN_SHINHAN     BankBIN = "970424" // Shinhan Bank
	BIN_STANCHART   BankBIN = "970410" // Standard Chartered
	BIN_TECHCOMBANK BankBIN = "970407" // Techcombank
	BIN_TIMO        BankBIN = "963388" // Timo
	BIN_TPBANK      BankBIN = "970423" // TPBank
	BIN_UOB         BankBIN = "970458" // United Overseas Bank
	BIN_VIB         BankBIN = "970441" // VIB
	BIN_VIETABANK   BankBIN = "970427" // VietA Bank
	BIN_VIETBANK    BankBIN = "970433" // VietBank
	BIN_VIETCOMBANK BankBIN = "970436" // Vietcombank
	BIN_VIETINBANK  BankBIN = "970415" // VietinBank
	BIN_VPBANK      BankBIN = "970432" // VPBank
	BIN_VRB         BankBIN = "970421" // Vietnam-Russia Joint Venture Bank
	BIN_WOORI       BankBIN = "970457" // Woori Bank
)

// BinForKey returns the Napas BIN for the given bank key, or empty if unknown.
func BinForKey(key BankKey) BankBIN {
	if b, ok := BanksMap[key]; ok {
		return BankBIN(b.BIN)
	}
	return ""
}
