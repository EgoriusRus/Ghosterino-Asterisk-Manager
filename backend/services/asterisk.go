package services

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"asterisk-manager/domain"
	"asterisk-manager/repositories"
)

// PhoneRecord представляет запись о телефоне/сотруднике
type PhoneRecord struct {
	FullName       string // A - ФИО
	Position       string // B - Должность
	CityPhone      string // C - № телефона (24-48-44)
	Extension      string // D - внутренний номер (4 цифры)
	PickupGroup    string // E - *8 pickupgroup/callgroup
	RingGroup      string // F - Группа входящих звонков
	IsT27          bool   // G - T27G
	IsT23          bool   // H - T23G
	IsRadio        bool   // I - Радио
	IsCiscoOrFax   bool   // J - Cisco или Fax
	IsActive       bool   // K - Включить в конфигурацию
	MACAddress     string // L - MAC телефона
	IsMobileClient bool   // M - Сотовый клиент
	VoipVLAN       string // N - voip vlan
	LanVLAN        string // O - lan vlan
	Location       string // P - Здание / адрес
	SIPServer      string // Q - ip сервера
	Subnet         string // R - ip подсеть телефона
	Email          string // S - email
	ConfRoom       string // T - конференц комнаты
	VoiceMenu      string // U - Дорожка голосового меню
	ExtraRingGroup string // V - доп ринг группа
	SpecialPass    string // W - Спец Пароль
	IsTLS          bool   // X - TLS
	IsFanvil       bool   // Y - fanvil
}

// CityNumber возвращает городской номер в формате 6 цифр (244844)
func (r *PhoneRecord) CityNumber() string {
	return strings.ReplaceAll(r.CityPhone, "-", "")
}

// GetPassword возвращает пароль (спец или MD5-хеш)
func (r *PhoneRecord) GetPassword() string {
	if r.SpecialPass != "" {
		return r.SpecialPass
	}
	return GetMD5Hash(r.Extension + r.Extension)
}

// IsLocalOnly проверяет, только локальные звонки
func (r *PhoneRecord) IsLocalOnly() bool {
	return r.RingGroup == "local"
}

// GetMD5Hash генерирует MD5-хеш строки
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// AsteriskGenerator генератор конфигурационных файлов
type AsteriskGenerator struct {
	Records   []PhoneRecord
	OutputDir string
}

// NewAsteriskGenerator создаёт новый генератор
func NewAsteriskGenerator(outputDir string) *AsteriskGenerator {
	return &AsteriskGenerator{
		OutputDir: outputDir,
		Records:   make([]PhoneRecord, 0),
	}
}

// LoadCSV загружает и парсит CSV файл
func (g *AsteriskGenerator) LoadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Разрешаем разное количество полей

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("ошибка чтения CSV: %w", err)
	}

	for i, row := range records {
		if i == 0 {
			continue // Пропускаем заголовки
		}

		if g.isServiceRow(row) {
			continue
		}

		record := g.parseRow(row)
		if record != nil {
			g.Records = append(g.Records, *record)
		}
	}

	fmt.Printf("Загружено записей: %d\n", len(g.Records))
	return nil
}

// LoadFromDatabase загружает данные из базы данных
func (g *AsteriskGenerator) LoadFromDatabase(repos *repositories.Repos) error {
	isActive := true
	profiles, _, err := repos.FindProfilesWithLocations(&isActive, nil)
	if err != nil {
		return fmt.Errorf("ошибка загрузки профилей: %w", err)
	}

	// Загружаем устройства
	var devices []domain.Device
	deviceMap := make(map[string]domain.Device)
	if err := repos.FindAllDevices(&devices); err != nil {
		return fmt.Errorf("ошибка загрузки устройств: %w", err)
	}
	for _, dev := range devices {
		deviceMap[dev.MAC] = dev
	}

	// Конвертируем domain-модели в PhoneRecord
	for _, profile := range profiles {
		record := profileToPhoneRecord(profile, deviceMap)
		if record != nil {
			g.Records = append(g.Records, *record)
		}
	}

	fmt.Printf("Загружено записей из БД: %d\n", len(g.Records))
	return nil
}

// profileToPhoneRecord конвертирует domain.ProfileWithLocation в PhoneRecord
func profileToPhoneRecord(p domain.ProfileWithLocation, deviceMap map[string]domain.Device) *PhoneRecord {
	if p.LocationName == nil || p.Server == nil || p.Subnet == nil {
		return nil // Пропускаем профили без локации
	}

	record := &PhoneRecord{
		FullName:  p.Name,
		Email:     p.Email,
		Extension: fmt.Sprintf("%d", p.InternalNumber),
		CityPhone: p.ExternalNumber,
		Location:  *p.LocationName,
		SIPServer: *p.Server,
		Subnet:    *p.Subnet,
		VoipVLAN:  strconv.Itoa(*p.VoipVLAN),
		LanVLAN:   strconv.Itoa(*p.VLAN),
		IsActive:  p.IsActive,
	}

	// Ring Group
	if p.RingGroup != nil {
		record.RingGroup = strconv.Itoa(*p.RingGroup)
	} else {
		record.RingGroup = "local"
	}

	// Pickup Group
	if p.PickupGroup != nil {
		record.PickupGroup = strconv.Itoa(*p.PickupGroup)
	}

	// Устройство
	if p.Device != nil {
		record.MACAddress = strings.ToLower(*p.Device)
		if dev, ok := deviceMap[*p.Device]; ok {
			switch dev.DeviceModel {
			case domain.DeviceModelYealinkT27G:
				record.IsT27 = true
			case domain.DeviceModelYealinkT23G:
				record.IsT23 = true
			case domain.DeviceModelFanvil:
				record.IsFanvil = true
			case domain.DeviceModelCisco:
				record.IsCiscoOrFax = true
			}
		}
	}

	return record
}

// isServiceRow проверяет, является ли строка служебной
func (g *AsteriskGenerator) isServiceRow(row []string) bool {
	if len(row) == 0 {
		return true
	}

	name := strings.TrimSpace(row[0])
	if name == "" {
		return true
	}

	// Служебные строки: "1,Надымская 3" или "2,Отдел..."
	// Первый символ - цифра, нет внутреннего номера
	if len(name) > 0 && name[0] >= '0' && name[0] <= '9' {
		// Проверяем внутренний номер (колонка D)
		if len(row) < 4 || strings.TrimSpace(row[3]) == "" {
			return true
		}
	}

	return false
}

// parseRow парсит строку CSV в PhoneRecord
func (g *AsteriskGenerator) parseRow(row []string) *PhoneRecord {
	if len(row) < 17 {
		return nil
	}

	extension := strings.TrimSpace(getField(row, 3))
	if extension == "" || len(extension) != 4 {
		return nil
	}

	record := &PhoneRecord{
		FullName:       strings.TrimSpace(getField(row, 0)),
		Position:       strings.TrimSpace(getField(row, 1)),
		CityPhone:      strings.TrimSpace(getField(row, 2)),
		Extension:      extension,
		PickupGroup:    strings.TrimSpace(getField(row, 4)),
		RingGroup:      strings.TrimSpace(getField(row, 5)),
		IsT27:          getField(row, 6) == "1",
		IsT23:          getField(row, 7) == "1",
		IsRadio:        getField(row, 8) == "1",
		IsCiscoOrFax:   getField(row, 9) == "1",
		IsActive:       getField(row, 10) == "1",
		MACAddress:     strings.ToLower(strings.TrimSpace(getField(row, 11))),
		IsMobileClient: getField(row, 12) == "1",
		VoipVLAN:       strings.TrimSpace(getField(row, 13)),
		LanVLAN:        strings.TrimSpace(getField(row, 14)),
		Location:       strings.TrimSpace(getField(row, 15)),
		SIPServer:      strings.TrimSpace(getField(row, 16)),
		Subnet:         strings.TrimSpace(getField(row, 17)),
		Email:          strings.TrimSpace(getField(row, 18)),
		ConfRoom:       strings.TrimSpace(getField(row, 19)),
		VoiceMenu:      strings.TrimSpace(getField(row, 20)),
		ExtraRingGroup: strings.TrimSpace(getField(row, 21)),
		SpecialPass:    strings.TrimSpace(getField(row, 22)),
		IsTLS:          getField(row, 23) == "1",
		IsFanvil:       getField(row, 24) == "1",
	}

	return record
}

func getField(row []string, index int) string {
	if index < len(row) {
		return row[index]
	}
	return ""
}

// Generate генерирует все конфигурационные файлы
func (g *AsteriskGenerator) Generate() error {
	// Создаём директории
	dirs := []string{
		filepath.Join(g.OutputDir, "tftpboot"),
		filepath.Join(g.OutputDir, "UsersConf"),
		filepath.Join(g.OutputDir, "ExtConf"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("ошибка создания директории %s: %w", dir, err)
		}
	}

	// Генерируем файлы
	if err := g.generateTftpboot(); err != nil {
		return err
	}

	if err := g.generateUsersConf(); err != nil {
		return err
	}

	if err := g.generateExtConf(); err != nil {
		return err
	}

	if err := g.generateCiscoConf(); err != nil {
		return err
	}

	return nil
}

// generateTftpboot генерирует конфиги для телефонов
func (g *AsteriskGenerator) generateTftpboot() error {
	dir := filepath.Join(g.OutputDir, "tftpboot")

	for _, r := range g.Records {
		// Условия из VBA:
		// - IsActive = true
		// - MAC >= 12 символов
		// - VoipVLAN, LanVLAN, SIPServer не пустые
		// - Тип телефона: T27 OR T23 OR Fanvil
		// - НЕ Cisco/Fax
		// - НЕ Сотовый клиент
		// - НЕ ExtraRingGroup
		if !r.IsActive {
			continue
		}
		if r.MACAddress == "" || len(r.MACAddress) < 12 {
			continue
		}
		if r.VoipVLAN == "" || r.LanVLAN == "" || r.SIPServer == "" {
			continue
		}
		if !r.IsT27 && !r.IsT23 && !r.IsFanvil {
			continue
		}
		if r.IsCiscoOrFax {
			continue
		}
		if r.IsMobileClient {
			continue
		}
		if r.ExtraRingGroup == "1" {
			continue
		}

		filename := filepath.Join(dir, r.MACAddress+".cfg")
		var content string

		if r.IsFanvil {
			content = g.generateFanvilConfig(r)
		} else if r.IsT23 {
			content = g.generateYealinkT23Config(r)
		} else if r.IsT27 {
			content = g.generateYealinkT27Config(r)
		}

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("ошибка записи %s: %w", filename, err)
		}
	}

	fmt.Println("✓ tftpboot конфиги сгенерированы")
	return nil
}

func (g *AsteriskGenerator) generateFanvilConfig(r PhoneRecord) string {
	var sb strings.Builder

	sb.WriteString("#<Voip Config File>#\n")
	sb.WriteString("Version = 2.0000000000\n")
	sb.WriteString(fmt.Sprintf("sip.line.1.PhoneNumber = %s\n", r.Extension))
	sb.WriteString(fmt.Sprintf("sip.line.1.DisplayName = %s\n", r.FullName))
	sb.WriteString(fmt.Sprintf("sip.line.1.SipName = %s\n", r.SIPServer))
	sb.WriteString(fmt.Sprintf("sip.line.1.RegAddr = %s\n", r.SIPServer))
	sb.WriteString(fmt.Sprintf("sip.line.1.RegUser = %s\n", r.Extension))
	sb.WriteString(fmt.Sprintf("sip.line.1.RegPswd = %s\n", r.GetPassword()))
	sb.WriteString("sip.line.1.RegEnabled = 1\n")
	sb.WriteString("sip.line.1.HotlineNum = 112\n")
	sb.WriteString("sip.line.1.HotlineEnabled = 0\n")
	sb.WriteString("sip.line.1.MWICode = *97\n")
	sb.WriteString("sip.line.1.AudioCodecSet = PCMU,PCMA,G726-32,G729,G722\n\n")

	// Сеть
	sb.WriteString("net.dhcp.Enabled = 1\n")
	sb.WriteString("net.dhcp.AutoDNS = 1\n")
	sb.WriteString("net.pppoe.Enabled = 0\n\n")

	// Телефон
	sb.WriteString("phone.MenuPassword = 123\n")
	sb.WriteString("phone.KeyLockPassword = 123\n")
	sb.WriteString("phone.KeyLockEnabled = 0\n")
	sb.WriteString("phone.KeyLockTimeout = 0\n")
	sb.WriteString("phone.KeyLockStatus = 0\n")
	sb.WriteString("phone.date.SNTPEnabled = 1\n")
	sb.WriteString("phone.date.SNTPServer = 10.16.0.100\n")
	sb.WriteString("phone.date.SecondSNTPServer = 10.16.0.69\n")
	sb.WriteString("phone.date.TimeZone = 20\n")
	sb.WriteString("phone.date.TimeZoneName = UTC+5\n")
	sb.WriteString("phone.date.SNTPInterval = 1000\n\n")

	// Веб
	sb.WriteString("web.account.1.Name = admin\n")
	sb.WriteString("web.account.1.Password = Admin19\n")
	sb.WriteString("web.account.1.Level = 10\n")
	sb.WriteString("web.account.2.Name = guest\n")
	sb.WriteString("web.account.2.Password = Guest19\n")
	sb.WriteString("web.account.2.Level = 5\n\n")

	// Provisioning
	sb.WriteString("ap.DefaultUsername = \n")
	sb.WriteString("ap.DefaultPassword = \n")
	sb.WriteString("ap.DownloadCommonConf = 1\n")
	sb.WriteString("ap.SaveProvisionInfo = 0\n")
	sb.WriteString("ap.FailedRetryTimes = 5\n")
	sb.WriteString("ap.FlashServerIP = 10.16.0.102\n")
	sb.WriteString("ap.FlashFileName = \n")
	sb.WriteString("ap.FlashProtocol = 2\n")
	sb.WriteString("ap.FlashMode = 1\n")
	sb.WriteString("ap.FlashInterval = 1\n")
	sb.WriteString("ap.DHCPOption = 66\n")
	sb.WriteString("ap.pnp.Enabled = 0\n")
	sb.WriteString("ap.pnp.IP = 10.16.0.102\n")
	sb.WriteString("ap.pnp.Port = 5060\n")
	sb.WriteString("ap.pnp.Transport = 0\n")
	sb.WriteString("ap.pnp.Interval = 1\n\n")

	// VLAN
	sb.WriteString("qos.VLANEnabled = 1\n")
	sb.WriteString(fmt.Sprintf("qos.VLANID = %s\n", r.VoipVLAN))
	sb.WriteString("qos.PortVLanEnabled = 1\n")
	sb.WriteString(fmt.Sprintf("qos.PortVLanID = %s\n", r.LanVLAN))

	return sb.String()
}

func (g *AsteriskGenerator) generateYealinkT23Config(r PhoneRecord) string {
	var sb strings.Builder

	sb.WriteString("#!version:1.0.0.1\n")
	sb.WriteString("#T23G\n")
	sb.WriteString(fmt.Sprintf("#%s\n", r.FullName))
	sb.WriteString("static.network.vlan.internet_port_enable = 1\n")
	sb.WriteString(fmt.Sprintf("static.network.vlan.internet_port_vid = %s\n", r.VoipVLAN))
	sb.WriteString("static.network.vlan.pc_port_enable = 0\n")
	sb.WriteString(fmt.Sprintf("static.network.vlan.pc_port_vid = %s\n", r.LanVLAN))
	sb.WriteString(fmt.Sprintf("account.1.auth_name = %s\n", r.Extension))
	sb.WriteString("account.1.codec.opus.enable = 0\n")
	sb.WriteString("account.1.codec.opus.priority = 5\n")
	sb.WriteString(fmt.Sprintf("account.1.display_name = %s %s\n", r.Extension, r.FullName))
	sb.WriteString("account.1.enable = 1\n")
	sb.WriteString(fmt.Sprintf("account.1.label = %s\n", r.Extension))
	sb.WriteString("account.1.ringtone.ring_type = Resource:Ring1.wav\n")
	sb.WriteString(fmt.Sprintf("account.1.sip_server.1.address = %s\n", r.SIPServer))
	sb.WriteString("account.1.subscribe_register = 1\n")
	sb.WriteString("account.1.unregister_on_reboot = 0\n")
	sb.WriteString(fmt.Sprintf("account.1.user_name = %s\n", r.Extension))
	sb.WriteString(fmt.Sprintf("static.network.dhcp_host_name = SIP-T23G-%s\n", r.Extension))
	sb.WriteString(fmt.Sprintf("account.1.password = %s\n", r.GetPassword()))

	return sb.String()
}

func (g *AsteriskGenerator) generateYealinkT27Config(r PhoneRecord) string {
	var sb strings.Builder

	sb.WriteString("#!version:1.0.0.1\n")
	sb.WriteString("#T27G\n")
	sb.WriteString("static.network.vlan.internet_port_enable = 1\n")
	sb.WriteString(fmt.Sprintf("static.network.vlan.internet_port_vid = %s\n", r.VoipVLAN))
	sb.WriteString("static.network.vlan.pc_port_enable = 0\n")
	sb.WriteString(fmt.Sprintf("static.network.vlan.pc_port_vid = %s\n", r.LanVLAN))
	sb.WriteString(fmt.Sprintf("account.1.auth_name = %s\n", r.Extension))
	sb.WriteString("account.1.codec.opus.enable = 0\n")
	sb.WriteString("account.1.codec.opus.priority = 5\n")
	sb.WriteString(fmt.Sprintf("account.1.display_name = %s %s\n", r.Extension, r.FullName))
	sb.WriteString("account.1.enable = 1\n")
	sb.WriteString(fmt.Sprintf("account.1.label = %s\n", r.Extension))
	sb.WriteString("account.1.ringtone.ring_type = Resource:Ring7.wav\n")
	sb.WriteString(fmt.Sprintf("account.1.sip_server.1.address = %s\n", r.SIPServer))
	sb.WriteString("account.1.subscribe_register = 1\n")
	sb.WriteString("account.1.unregister_on_reboot = 0\n")
	sb.WriteString(fmt.Sprintf("static.network.dhcp_host_name = SIP-T27G-%s\n", r.Extension))
	sb.WriteString(fmt.Sprintf("account.1.user_name = %s\n", r.Extension))
	sb.WriteString(fmt.Sprintf("account.1.password = %s\n", r.GetPassword()))
	sb.WriteString("distinctive_ring_tones.alert_info.1.ringer = Resource:Ring2.wav\n")
	sb.WriteString("linekey.9.xml_phonebook = 0\n")

	return sb.String()
}

// HasDeviceType проверяет наличие типа устройства
func (r *PhoneRecord) HasDeviceType() bool {
	return r.IsT27 || r.IsT23 || r.IsRadio || r.IsCiscoOrFax || r.IsMobileClient || r.IsFanvil
}

// generateUsersConf генерирует конфиги пользователей SIP
func (g *AsteriskGenerator) generateUsersConf() error {
	dir := filepath.Join(g.OutputDir, "UsersConf")

	for _, r := range g.Records {
		// Условия из VBA:
		// - IsActive = true
		// - Location и Subnet не пустые
		// - Есть тип устройства (T27/T23/Радио/Cisco/Сотовый/Fanvil)
		// - НЕ ExtraRingGroup
		if !r.IsActive {
			continue
		}
		if r.Location == "" || r.Subnet == "" {
			continue
		}
		if !r.HasDeviceType() {
			continue
		}
		if r.ExtraRingGroup == "1" {
			continue
		}

		filename := filepath.Join(dir, fmt.Sprintf("User%s.conf", r.Extension))
		content := g.generateUserConfig(r)

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("ошибка записи %s: %w", filename, err)
		}
	}

	fmt.Println("✓ UsersConf конфиги сгенерированы")
	return nil
}

func (g *AsteriskGenerator) generateUserConfig(r PhoneRecord) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s]\n", r.Extension))
	sb.WriteString(fmt.Sprintf("fullname = %s\n", r.FullName))
	sb.WriteString("registersip = no\n")
	sb.WriteString("host = Dynamic\n")

	if r.PickupGroup != "" {
		sb.WriteString(fmt.Sprintf("callgroup = %s\n\n", r.PickupGroup))
	} else {
		sb.WriteString(";callgroup = \n\n")
	}

	sb.WriteString(fmt.Sprintf("mailbox = %s\n", r.Extension))
	sb.WriteString("call-limit = 100\n")
	sb.WriteString("type = peer\n")
	sb.WriteString(fmt.Sprintf("UserName = %s\n", r.Extension))
	sb.WriteString("transfer = yes\n\n")

	sb.WriteString("callcounter = yes\n")

	// Context
	if r.IsLocalOnly() {
		sb.WriteString("Context = DLPN_DialPlan_OnlyLocal\n")
	} else {
		sb.WriteString(fmt.Sprintf("Context = DLPN_DialPlan_%s_%s\n", r.Location, r.CityNumber()))
	}

	sb.WriteString(fmt.Sprintf("cid_number = %s\n", r.Extension))
	sb.WriteString("hasvoicemail = yes\n")
	sb.WriteString("vmsecret = 1234\n")
	sb.WriteString(fmt.Sprintf("email = %s\n", r.Email))
	sb.WriteString("threewaycalling = no\n")
	sb.WriteString("hasdirectory = no\n")
	sb.WriteString("callwaiting = no\n")
	sb.WriteString("hasmanager = no\n")
	sb.WriteString("hasagent = no\n")
	sb.WriteString("hassip = yes\n")
	sb.WriteString("hasiax = no\n")
	sb.WriteString(fmt.Sprintf("secret = %s\n", r.GetPassword()))
	sb.WriteString("nat=force_rport,comedia\n")
	sb.WriteString("canreinvite = no\n")
	sb.WriteString("dtmfmode = rfc2833\n")
	sb.WriteString("insecure = no\n")

	if r.PickupGroup != "" {
		sb.WriteString(fmt.Sprintf("pickupgroup = %s\n\n", r.PickupGroup))
	} else {
		sb.WriteString(";pickupgroup = \n\n")
	}

	sb.WriteString(fmt.Sprintf("macaddress = %s\n", r.Extension))
	sb.WriteString("autoprov = yes\n")
	sb.WriteString(fmt.Sprintf("Label = %s\n", r.Extension))
	sb.WriteString("linenumber = 1\n")
	sb.WriteString("LINEKEYS = 1\n")
	sb.WriteString("disallow = all\n")
	sb.WriteString("allow = alaw,ulaw\n")
	sb.WriteString("deny = 0.0.0.0/0\n")
	sb.WriteString(fmt.Sprintf("permit = %s\n", r.Subnet))
	sb.WriteString("directmedia = no\n")

	return sb.String()
}

// generateExtConf генерирует файлы диалплана
func (g *AsteriskGenerator) generateExtConf() error {
	dir := filepath.Join(g.OutputDir, "ExtConf")

	// ExtensionsCID.conf
	if err := g.generateExtensionsCID(dir); err != nil {
		return err
	}

	// ExtensionsRG.conf и ExtensionsRGCFG.conf
	if err := g.generateRingGroups(dir); err != nil {
		return err
	}

	// ExtensionsOut.conf и ExtensionsDP.conf
	if err := g.generateDialplans(dir); err != nil {
		return err
	}

	// ExtensionsTrunkZags.conf и ExtensionsTrankAdm.conf
	if err := g.generateTrunks(dir); err != nil {
		return err
	}

	fmt.Println("✓ ExtConf конфиги сгенерированы")
	return nil
}

func (g *AsteriskGenerator) generateExtensionsCID(dir string) error {
	var sb strings.Builder
	seen := make(map[string]bool)

	for _, r := range g.Records {
		if !r.IsActive || seen[r.Extension] {
			continue
		}
		seen[r.Extension] = true
		sb.WriteString(fmt.Sprintf("CID_%s = %s\n", r.Extension, r.Extension))
	}

	return os.WriteFile(filepath.Join(dir, "ExtensionsCID.conf"), []byte(sb.String()), 0644)
}

func (g *AsteriskGenerator) generateRingGroups(dir string) error {
	// Собираем ринг-группы
	ringGroups := make(map[string][]PhoneRecord)

	for _, r := range g.Records {
		if !r.IsActive || r.IsLocalOnly() || r.RingGroup == "" {
			continue
		}
		ringGroups[r.RingGroup] = append(ringGroups[r.RingGroup], r)
	}

	// ExtensionsRG.conf
	var sbRG strings.Builder
	sbRG.WriteString("[ringgroups]\n")

	// ExtensionsRGCFG.conf
	var sbRGCFG strings.Builder
	sbRGCFG.WriteString(";Настройки ринг групп\n")

	// ExtensionsVMCFG.conf
	var sbVMCFG strings.Builder
	sbVMCFG.WriteString(";Голосовое меню\n")

	for rgNum, members := range ringGroups {
		if len(members) == 0 {
			continue
		}

		first := members[0]
		cityNum := first.CityNumber()

		// RG entry
		sbRG.WriteString(fmt.Sprintf("exten => %s,1,Goto(ringroups-%s-%s,s,1)\n", rgNum, cityNum, rgNum))

		// RGCFG
		sbRGCFG.WriteString(fmt.Sprintf("[ringroups-%s-%s]\n", cityNum, rgNum))
		sbRGCFG.WriteString(fmt.Sprintf("exten => s,1,NoOp(RG%s)\n", cityNum))

		// Dial string
		var dialParts []string
		for _, m := range members {
			dialParts = append(dialParts, fmt.Sprintf("SIP/%s", m.Extension))
		}
		sbRGCFG.WriteString(fmt.Sprintf("exten => s,n,Dial(%s,40,${DIALOPTIONS})\n", strings.Join(dialParts, "&")))
		sbRGCFG.WriteString(fmt.Sprintf("exten => s,n,Voicemail(%s,u)\n\n", first.Extension))

		// VMCFG
		sbVMCFG.WriteString(fmt.Sprintf("[voicemenu-%s-%s]\n", cityNum, rgNum))
		sbVMCFG.WriteString(fmt.Sprintf("exten => s,1,NoOp(VM%s)\n", cityNum))
		sbVMCFG.WriteString("exten => s,n,Set(numTries=0)\n")
		sbVMCFG.WriteString("exten => s,n,Answer()\n")

		if first.VoiceMenu != "" && first.VoiceMenu != "NO" {
			sbVMCFG.WriteString(fmt.Sprintf("exten => s,n(naberite),Background(record/%s)\n", first.VoiceMenu))
			sbVMCFG.WriteString("exten => s,n,WaitExten(5)\n")
		}

		sbVMCFG.WriteString("exten => s,n,Background(record/WRITECALL)\n")
		sbVMCFG.WriteString(fmt.Sprintf("exten => s,n,Goto(ringroups-%s-%s,s,1)\n", cityNum, rgNum))

		// Прямые номера в меню
		for _, m := range members {
			sbVMCFG.WriteString(fmt.Sprintf("exten => _%s,1,Background(record/WRITECALL)\n", m.Extension))
			sbVMCFG.WriteString(fmt.Sprintf("exten => _%s,2,Dial(SIP/${EXTEN})\n", m.Extension))
		}

		sbVMCFG.WriteString("exten => _XXXX,1,Goto(s,naberite)\n")
		sbVMCFG.WriteString("exten => 0,1,Goto(s,naberite)\n")

		// Факс
		sbVMCFG.WriteString(fmt.Sprintf("exten => 1,1,Set(FAXFILE=/var/calls/FAX/${STRFTIME(${EPOCH},,%s-%%Y%%m%%d-%%H_%%M_%%S)}-from-${CALLERID(num)}.tif)\n", cityNum))
		sbVMCFG.WriteString(fmt.Sprintf("exten => 1,2,Set(PDFFILE=/var/calls/FAX/${STRFTIME(${EPOCH},,%s-%%Y%%m%%d-%%H_%%M_%%S)}-from-${CALLERID(num)}.pdf)\n", cityNum))
		sbVMCFG.WriteString("exten => 1,3,ReceiveFax(${FAXFILE})\n")
		sbVMCFG.WriteString("exten => 1,4,System(/usr/bin/tiff2pdf ${FAXFILE} > ${PDFFILE})\n")
		sbVMCFG.WriteString("exten => 1,5,System(/bin/rm -f ${FAXFILE})\n")

		email := first.Email
		if email == "" {
			email = "fax@nur.yanao.ru"
		}
		sbVMCFG.WriteString(fmt.Sprintf("exten => 1,6,System(/root/bin/sendEmail.pl -f fax%s@nur.yanao.ru -t %s -u \"Incoming FAX ${CALLERID(num)}\" -m \"Вам пришел факс с номера ${CALLERID(num)} в ${STRFTIME(${EPOCH},,%%H:%%M:%%S)}. Факс во вложении.\" -a ${PDFFILE} -o message-charset=UTF-8)\n", cityNum, email))
		sbVMCFG.WriteString("exten => 1,7,Hangup()\n")
		sbVMCFG.WriteString("exten => 2,1,Background(record/VoiceMesAns)\n")
		sbVMCFG.WriteString(fmt.Sprintf("exten => 2,2,Voicemail(%s,s)\n\n", first.Extension))
	}

	// Статические записи для ExtensionsRG.conf (из VBA)
	sbRG.WriteString("\n; Статические записи\n")
	sbRG.WriteString("exten => 6246,1,Goto(ringroups-947947-6246,s,1)\n")
	sbRG.WriteString("exten => 6293,1,Goto(ringroups-947994-6293,s,1)\n")
	sbRG.WriteString("exten => 6097,1,Goto(ringroups-947798-6097,s,1)\n")

	if err := os.WriteFile(filepath.Join(dir, "ExtensionsRG.conf"), []byte(sbRG.String()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "ExtensionsRGCFG.conf"), []byte(sbRGCFG.String()), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "ExtensionsVMCFG.conf"), []byte(sbVMCFG.String()), 0644)
}

func (g *AsteriskGenerator) generateDialplans(dir string) error {
	var sbOut strings.Builder
	var sbDP strings.Builder

	sbOut.WriteString(";Исходящие правила\n")
	sbDP.WriteString(";Перечень диалпланов\n")

	// OnlyLocal dialplan
	sbDP.WriteString("[DLPN_DialPlan_OnlyLocal]\n")
	sbDP.WriteString("include => CallingRule_DIOout\n")
	sbDP.WriteString("include => default\n")
	sbDP.WriteString("include => parkedcalls\n")
	sbDP.WriteString("include => conferences\n")
	sbDP.WriteString("include => ringgroups\n")
	sbDP.WriteString("include => voicemenus\n")
	sbDP.WriteString("include => queues\n")
	sbDP.WriteString("include => voicemailgroups\n")
	sbDP.WriteString("include => directory\n")
	sbDP.WriteString("include => pagegroups\n")
	sbDP.WriteString("include => page_an_extension\n\n")

	// Уникальные городские номера
	seenCityNum := make(map[string]bool)
	admLocations := map[string]bool{
		"Mir": true, "Sov": true, "Ubil1": true, "Leb5b-arh": true,
		"Nad3": true, "Limb": true, "Kor": true, "Ind4": true, "Len15v": true,
	}

	for _, r := range g.Records {
		if !r.IsActive || r.IsLocalOnly() {
			continue
		}

		cityNum := r.CityNumber()
		if cityNum == "" || len(cityNum) != 6 || seenCityNum[cityNum] {
			continue
		}
		seenCityNum[cityNum] = true

		trunk := "trunk_2"
		if admLocations[r.Location] {
			trunk = "trunk_3"
		}

		// Spec rules
		sbOut.WriteString(fmt.Sprintf("[CallingRule_RT%sout-spec]\n", cityNum))
		sbOut.WriteString("exten => _X.,1,Macro(recording,${CALLERID(num)},${EXTEN})\n")
		sbOut.WriteString(fmt.Sprintf("exten => _0X,2,Macro(trunkdial-failover-0.3,${%s}/${EXTEN:0},,%s,,3494%s)\n", trunk, trunk, cityNum))
		sbOut.WriteString(fmt.Sprintf("exten => _1XX,3,Macro(trunkdial-failover-0.3,${%s}/${EXTEN:0},,%s,,3494%s)\n\n", trunk, trunk, cityNum))

		// Main rules
		sbOut.WriteString(fmt.Sprintf("[CallingRule_RT%sout]\n", cityNum))
		sbOut.WriteString("exten => _X.,1,Macro(recording,${CALLERID(num)},${EXTEN})\n")
		sbOut.WriteString(fmt.Sprintf("exten => _[29]XXXXX,2,Macro(trunkdial-failover-0.3,${%s}/${EXTEN:0},,%s,,3494%s)\n", trunk, trunk, cityNum))
		sbOut.WriteString(fmt.Sprintf("exten => _[78]XXXXX.,3,Macro(trunkdial-failover-0.3,${%s}/8${EXTEN:1},,%s,,3494%s)\n", trunk, trunk, cityNum))
		sbOut.WriteString(fmt.Sprintf("exten => _NXXXXXXXXX,4,Macro(trunkdial-failover-0.3,${%s}/8${EXTEN:0},,%s,,3494%s)\n", trunk, trunk, cityNum))

		if r.ConfRoom != "" {
			sbOut.WriteString(fmt.Sprintf("exten => %s,1,Answer()\n", r.ConfRoom))
			sbOut.WriteString(fmt.Sprintf("exten => %s,n,ConfBridge(1,confer)\n", r.ConfRoom))
		}
		sbOut.WriteString("\n")

		// Dialplan
		sbDP.WriteString(fmt.Sprintf("[DLPN_DialPlan_%s_%s]\n", r.Location, cityNum))
		sbDP.WriteString(fmt.Sprintf("include => CallingRule_RT%sout-spec\n", cityNum))
		sbDP.WriteString("include => CallingRule_DIOout\n")
		sbDP.WriteString(fmt.Sprintf("include => CallingRule_RT%sout\n", cityNum))
		sbDP.WriteString("include => default\n")
		sbDP.WriteString("include => parkedcalls\n")
		sbDP.WriteString("include => conferences\n")
		sbDP.WriteString("include => ringgroups\n")
		sbDP.WriteString("include => voicemenus\n")
		sbDP.WriteString("include => queues\n")
		sbDP.WriteString("include => voicemailgroups\n")
		sbDP.WriteString("include => directory\n")
		sbDP.WriteString("include => pagegroups\n")
		sbDP.WriteString("include => page_an_extension\n\n")
	}

	if err := os.WriteFile(filepath.Join(dir, "ExtensionsOut.conf"), []byte(sbOut.String()), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "ExtensionsDP.conf"), []byte(sbDP.String()), 0644)
}

func (g *AsteriskGenerator) generateTrunks(dir string) error {
	var sbZags strings.Builder
	var sbAdm strings.Builder

	sbZags.WriteString(";Входящие - шлюз ЗАГС\n")
	sbZags.WriteString("[DID_trunk_2]\n")
	sbZags.WriteString("include => DID_trunk_2_default\n")
	sbZags.WriteString("[DID_trunk_2_default]\n")

	sbAdm.WriteString(";Входящие - шлюз Администрация\n")
	sbAdm.WriteString("[DID_trunk_3]\n")
	sbAdm.WriteString("include => DID_trunk_3_default\n")
	sbAdm.WriteString("[DID_trunk_3_default]\n")

	seenCityNum := make(map[string]bool)
	admLocations := map[string]bool{
		"Mir": true, "Limb": true, "Kor": true, "Sov": true,
		"Leb5b-arh": true, "Ubil1": true, "Ind4": true, "Len15v": true, "Nad3": true,
	}

	for _, r := range g.Records {
		if !r.IsActive || r.IsLocalOnly() {
			continue
		}

		cityNum := r.CityNumber()
		if cityNum == "" || len(cityNum) != 6 || seenCityNum[cityNum] {
			continue
		}
		seenCityNum[cityNum] = true

		entry := fmt.Sprintf("exten => _%s,1,Macro(recording,${CALLERID(num)},${EXTEN})\n", cityNum)
		entry += fmt.Sprintf("exten => _%s,2,Goto(voicemenu-%s-%s,s,1)\n", cityNum, cityNum, r.RingGroup)

		if r.Location == "Zags" {
			sbZags.WriteString(entry)
		} else if admLocations[r.Location] {
			sbAdm.WriteString(entry)
		}
	}

	// Статические записи для trunk_3 (Администрация) - из VBA
	sbAdm.WriteString("\n; Статические записи\n")
	sbAdm.WriteString("exten => _947947,1,Macro(recording,${CALLERID(num)},${EXTEN})\n")
	sbAdm.WriteString("exten => _947947,2,Goto(voicemenu-947947-6246,s,1)\n")
	sbAdm.WriteString("exten => _947994,1,Macro(recording,${CALLERID(num)},${EXTEN})\n")
	sbAdm.WriteString("exten => _947994,2,Goto(voicemenu-947994-6293,s,1)\n")
	sbAdm.WriteString("exten => _947798,1,Macro(recording,${CALLERID(num)},${EXTEN})\n")
	sbAdm.WriteString("exten => _947798,2,Goto(voicemenu-947798-6097,s,1)\n")

	if err := os.WriteFile(filepath.Join(dir, "ExtensionsTrunkZags.conf"), []byte(sbZags.String()), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "ExtensionsTrankAdm.conf"), []byte(sbAdm.String()), 0644)
}

// generateCiscoConf генерирует конфиг для Cisco
func (g *AsteriskGenerator) generateCiscoConf() error {
	var sb strings.Builder
	seenCityNum := make(map[string]bool)

	for _, r := range g.Records {
		if !r.IsActive || r.Location == "Zags" {
			continue
		}

		cityNum := r.CityNumber()
		if cityNum == "" || len(cityNum) != 6 || seenCityNum[cityNum] {
			continue
		}
		seenCityNum[cityNum] = true

		sb.WriteString(fmt.Sprintf("\ndial-peer voice %s voip\n", cityNum))
		sb.WriteString("corlist incoming pbx94\n")
		sb.WriteString(fmt.Sprintf("Description %s - %s\n", cityNum, r.Location))
		sb.WriteString("huntstop\n")
		sb.WriteString(fmt.Sprintf("destination-pattern %s\n", cityNum))
		sb.WriteString("voice-class codec 1\n")
		sb.WriteString("session protocol sipv2\n")
		sb.WriteString(fmt.Sprintf("session target ipv4:%s:5060\n", r.SIPServer))
		sb.WriteString("session transport udp\n")
		sb.WriteString("dtmf-relay rtp-nte\n")
		sb.WriteString("fax rate disable\n")
		sb.WriteString("fax protocol pass-through g711ulaw\n")
		sb.WriteString("no vad\n")
	}

	filename := filepath.Join(g.OutputDir, "CiscoConf.txt")
	if err := os.WriteFile(filename, []byte(sb.String()), 0644); err != nil {
		return err
	}

	fmt.Println("✓ CiscoConf.txt сгенерирован")
	return nil
}

// GetStats возвращает статистику
func (g *AsteriskGenerator) GetStats() map[string]int {
	stats := map[string]int{
		"total":   len(g.Records),
		"active":  0,
		"t27":     0,
		"t23":     0,
		"fanvil":  0,
		"cisco":   0,
		"withMAC": 0,
		"local":   0,
	}

	for _, r := range g.Records {
		if r.IsActive {
			stats["active"]++
		}
		if r.IsT27 {
			stats["t27"]++
		}
		if r.IsT23 {
			stats["t23"]++
		}
		if r.IsFanvil {
			stats["fanvil"]++
		}
		if r.IsCiscoOrFax {
			stats["cisco"]++
		}
		if r.MACAddress != "" {
			stats["withMAC"]++
		}
		if r.IsLocalOnly() {
			stats["local"]++
		}
	}

	return stats
}
