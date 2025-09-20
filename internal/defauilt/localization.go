package defauilt

import "strings"

// dictionary описывает набор локализованных строк.
type dictionary struct {
	StatusSuccess             string
	StatusProblem             string
	StatusInProgress          string
	SummaryCompleted          string
	SummaryOf                 string
	SummaryTasks              string
	DefaultNo                 string
	DefaultYes                string
	DefaultSuccessLabel       string
	DefaultFromSummaryLabel   string
	DefaultTasksSummaryLabel  string
	DefaultErrorLabel         string
	DefaultCancelLabel        string
	DefaultSelectedLabel      string
	DefaultYesLabel           string
	DefaultNoLabel            string
	TaskCancelledByUser       string
	TaskExitHint              string
	DefaultPrompt             string
	PasswordMask              rune
	DefaultPlaceholder        string
	DefaultSeparator          string
	ErrorTypeValidation       string
	ErrorTypeUserCancel       string
	ErrorTypeTimeout          string
	ErrorTypeNetwork          string
	ErrorTypeFileSystem       string
	ErrorTypePermission       string
	ErrorTypeConfig           string
	ErrorTypeUnknown          string
	ErrorMsgUnknown           string
	ErrorMsgTaskPrefix        string
	ErrorMsgCanceled          string
	ErrorMsgTimeout           string
	ErrorMsgPermission        string
	ErrorUserMsgValidation    string
	ErrorUserMsgCancel        string
	ErrorUserMsgTimeout       string
	ErrorUserMsgNetwork       string
	ErrorUserMsgFileSystem    string
	ErrorUserMsgPermission    string
	ErrorUserMsgConfiguration string
	ErrorUserMsgUnknown       string
	TaskStatusError           string
	TaskStatusCancelled       string
	ErrFieldRequired          string
	ErrPathEmpty              string
	ErrPathInvalidChar        string
	ErrURLEmpty               string
	ErrURLScheme              string
	ErrValueEmpty             string
	ErrValueAlphaNumeric      string
	ErrDefaultValueInvalid    string
	ErrDefaultValueEmpty      string
	CancelShort               string
	NeedSelectAtLeastOne      string
	ScrollAboveFormat         string
	ScrollBelowFormat         string
	SingleSelectHelp          string
	MultiSelectHelp           string
	MultiSelectHelpSelectAll  string
	SelectAllDefaultText      string
	InputConfirmHint          string
	InputFormatLabel          string
	InputHintPassword         string
	InputHintEmail            string
	InputHintNumber           string
	InputHintIP               string
	InputHintDomain           string
}

var (
	defaultLanguage = "ru"
	currentLanguage = "ru"

	dictionaries = map[string]dictionary{
		"ru": {
			StatusSuccess:             "УСПЕШНО",
			StatusProblem:             "ПРОБЛЕМА",
			StatusInProgress:          "В ПРОЦЕССЕ",
			SummaryCompleted:          "Успешно завершено",
			SummaryOf:                 "из",
			SummaryTasks:              "задач",
			DefaultNo:                 "Нет",
			DefaultYes:                "Да",
			DefaultSuccessLabel:       "Готово",
			DefaultFromSummaryLabel:   "из",
			DefaultTasksSummaryLabel:  "задач",
			DefaultErrorLabel:         "Ошибка",
			DefaultCancelLabel:        "Отменено пользователем",
			DefaultSelectedLabel:      "пользователь выбрал",
			DefaultYesLabel:           "УСПЕШНО",
			DefaultNoLabel:            "ОТКАЗ",
			TaskCancelledByUser:       "[отменено пользователем]",
			TaskExitHint:              "[Для выхода из задачи нажмите Ctrl+C]",
			DefaultPrompt:             "Введите значение",
			PasswordMask:              '*',
			DefaultPlaceholder:        "...",
			DefaultSeparator:          "♀ ",
			ErrorTypeValidation:       "ВАЛИДАЦИЯ",
			ErrorTypeUserCancel:       "ОТМЕНА",
			ErrorTypeTimeout:          "ТАЙМАУТ",
			ErrorTypeNetwork:          "СЕТЬ",
			ErrorTypeFileSystem:       "ФАЙЛ",
			ErrorTypePermission:       "ДОСТУП",
			ErrorTypeConfig:           "КОНФИГ",
			ErrorTypeUnknown:          "ОШИБКА",
			ErrorMsgUnknown:           "неизвестная ошибка",
			ErrorMsgTaskPrefix:        "задача '%s': ",
			ErrorMsgCanceled:          "отменено пользователем",
			ErrorMsgTimeout:           "операция не завершилась за %v",
			ErrorMsgPermission:        "недостаточно прав для доступа к %s",
			ErrorUserMsgValidation:    "Проверьте правильность введенных данных",
			ErrorUserMsgCancel:        "Операция отменена",
			ErrorUserMsgTimeout:       "Операция заняла слишком много времени",
			ErrorUserMsgNetwork:       "Проблема с сетевым соединением",
			ErrorUserMsgFileSystem:    "Проблема доступа к файлу",
			ErrorUserMsgPermission:    "Недостаточно прав для выполнения операции",
			ErrorUserMsgConfiguration: "Ошибка в настройках",
			ErrorUserMsgUnknown:       "Произошла неизвестная ошибка",
			TaskStatusError:           "ОШИБКА",
			TaskStatusCancelled:       "Отменено",
			ErrFieldRequired:          "поле обязательно для заполнения",
			ErrPathEmpty:              "путь не может быть пустым",
			ErrPathInvalidChar:        "путь содержит недопустимый символ: %c",
			ErrURLEmpty:               "URL не может быть пустым",
			ErrURLScheme:              "URL должен начинаться с http:// или https://",
			ErrValueEmpty:             "значение не может быть пустым",
			ErrValueAlphaNumeric:      "значение должно содержать только буквы и цифры",
			ErrDefaultValueInvalid:    "значение по умолчанию невалидно",
			ErrDefaultValueEmpty:      "значение по умолчанию пусто",
			CancelShort:               "Отменено",
			NeedSelectAtLeastOne:      "! Необходимо выбрать хотя бы один элемент",
			ScrollAboveFormat:         "%s %s %d выше",
			ScrollBelowFormat:         "%s %s %d ниже",
			SingleSelectHelp:          "[← выход, ↑/↓ навигация, →/Enter выбор, Q/Esc - выход]",
			MultiSelectHelp:           "[← выход, ↑/↓ навигация, →/пробел выбор, Enter подтверждение, Q/Esc - выход]",
			MultiSelectHelpSelectAll:  "[← выход, ↑/↓ навигация, →/пробел выбор/переключение всех, Enter подтверждение, Q/Esc - выход]",
			SelectAllDefaultText:      "Выбрать все",
			InputConfirmHint:          "[Enter - подтвердить, Ctrl+C - отменить]",
			InputFormatLabel:          "Формат:",
			InputHintPassword:         "Используйте надежный пароль",
			InputHintEmail:            "Пример: user@example.com",
			InputHintNumber:           "Введите число",
			InputHintIP:               "Пример: 192.168.1.1",
			InputHintDomain:           "Пример: example.com",
		},
		"en": {
			StatusSuccess:             "SUCCESS",
			StatusProblem:             "ISSUE",
			StatusInProgress:          "IN PROGRESS",
			SummaryCompleted:          "Completed",
			SummaryOf:                 "of",
			SummaryTasks:              "tasks",
			DefaultNo:                 "No",
			DefaultYes:                "Yes",
			DefaultSuccessLabel:       "Done",
			DefaultFromSummaryLabel:   "of",
			DefaultTasksSummaryLabel:  "tasks",
			DefaultErrorLabel:         "Error",
			DefaultCancelLabel:        "Cancelled by user",
			DefaultSelectedLabel:      "user selected",
			DefaultYesLabel:           "SUCCESS",
			DefaultNoLabel:            "DECLINED",
			TaskCancelledByUser:       "[cancelled by user]",
			TaskExitHint:              "[Press Ctrl+C to exit task]",
			DefaultPrompt:             "Enter value",
			PasswordMask:              '*',
			DefaultPlaceholder:        "...",
			DefaultSeparator:          ", ",
			ErrorTypeValidation:       "VALIDATION",
			ErrorTypeUserCancel:       "CANCEL",
			ErrorTypeTimeout:          "TIMEOUT",
			ErrorTypeNetwork:          "NETWORK",
			ErrorTypeFileSystem:       "FILESYSTEM",
			ErrorTypePermission:       "PERMISSION",
			ErrorTypeConfig:           "CONFIG",
			ErrorTypeUnknown:          "ERROR",
			ErrorMsgUnknown:           "unknown error",
			ErrorMsgTaskPrefix:        "task '%s': ",
			ErrorMsgCanceled:          "cancelled by user",
			ErrorMsgTimeout:           "operation timed out after %v",
			ErrorMsgPermission:        "insufficient permissions to access %s",
			ErrorUserMsgValidation:    "Check the entered data",
			ErrorUserMsgCancel:        "Operation cancelled",
			ErrorUserMsgTimeout:       "Operation took too long",
			ErrorUserMsgNetwork:       "Network issue detected",
			ErrorUserMsgFileSystem:    "Filesystem access problem",
			ErrorUserMsgPermission:    "Not enough privileges to complete the operation",
			ErrorUserMsgConfiguration: "Configuration issue detected",
			ErrorUserMsgUnknown:       "An unknown error occurred",
			TaskStatusError:           "ERROR",
			TaskStatusCancelled:       "Cancelled",
			ErrFieldRequired:          "field is required",
			ErrPathEmpty:              "path cannot be empty",
			ErrPathInvalidChar:        "path contains an invalid character: %c",
			ErrURLEmpty:               "URL cannot be empty",
			ErrURLScheme:              "URL must start with http:// or https://",
			ErrValueEmpty:             "value cannot be empty",
			ErrValueAlphaNumeric:      "value must contain only letters and digits",
			ErrDefaultValueInvalid:    "default value is invalid",
			ErrDefaultValueEmpty:      "default value is empty",
			CancelShort:               "Cancelled",
			NeedSelectAtLeastOne:      "! Select at least one item",
			ScrollAboveFormat:         "%s %s %d above",
			ScrollBelowFormat:         "%s %s %d below",
			SingleSelectHelp:          "[← exit, ↑/↓ navigate, →/Enter select, Q/Esc exit]",
			MultiSelectHelp:           "[← exit, ↑/↓ navigate, →/space toggle, Enter confirm, Q/Esc exit]",
			MultiSelectHelpSelectAll:  "[← exit, ↑/↓ navigate, →/space toggle all, Enter confirm, Q/Esc exit]",
			SelectAllDefaultText:      "Select all",
			InputConfirmHint:          "[Enter to confirm, Ctrl+C to cancel]",
			InputFormatLabel:          "Format:",
			InputHintPassword:         "Use a reliable password",
			InputHintEmail:            "Example: user@example.com",
			InputHintNumber:           "Enter a number",
			InputHintIP:               "Example: 192.168.1.1",
			InputHintDomain:           "Example: example.com",
		},
		"tr": {
			StatusSuccess:             "BAŞARILI",
			StatusProblem:             "SORUN",
			StatusInProgress:          "DEVAM EDİYOR",
			SummaryCompleted:          "Başarıyla tamamlandı",
			SummaryOf:                 "/",
			SummaryTasks:              "görev",
			DefaultNo:                 "Hayır",
			DefaultYes:                "Evet",
			DefaultSuccessLabel:       "Tamamlandı",
			DefaultFromSummaryLabel:   "/",
			DefaultTasksSummaryLabel:  "görev",
			DefaultErrorLabel:         "Hata",
			DefaultCancelLabel:        "Kullanıcı iptal etti",
			DefaultSelectedLabel:      "kullanıcı seçti",
			DefaultYesLabel:           "BAŞARILI",
			DefaultNoLabel:            "REDDEDİLDİ",
			TaskCancelledByUser:       "[kullanıcı iptal etti]",
			TaskExitHint:              "[Görevden çıkmak için Ctrl+C]",
			DefaultPrompt:             "Değer girin",
			PasswordMask:              '*',
			DefaultPlaceholder:        "...",
			DefaultSeparator:          ", ",
			ErrorTypeValidation:       "DOĞRULAMA",
			ErrorTypeUserCancel:       "İPTAL",
			ErrorTypeTimeout:          "ZAMAN AŞIMI",
			ErrorTypeNetwork:          "AĞ",
			ErrorTypeFileSystem:       "DOSYA",
			ErrorTypePermission:       "İZİN",
			ErrorTypeConfig:           "KONFIG",
			ErrorTypeUnknown:          "HATA",
			ErrorMsgUnknown:           "bilinmeyen hata",
			ErrorMsgTaskPrefix:        "görev '%s': ",
			ErrorMsgCanceled:          "kullanıcı iptal etti",
			ErrorMsgTimeout:           "işlem %v sürede tamamlanamadı",
			ErrorMsgPermission:        "%s için yeterli izin yok",
			ErrorUserMsgValidation:    "Girilen verileri kontrol edin",
			ErrorUserMsgCancel:        "İşlem iptal edildi",
			ErrorUserMsgTimeout:       "İşlem çok uzun sürdü",
			ErrorUserMsgNetwork:       "Ağ bağlantı sorunu",
			ErrorUserMsgFileSystem:    "Dosya sistemi erişim sorunu",
			ErrorUserMsgPermission:    "İşlem için yeterli yetki yok",
			ErrorUserMsgConfiguration: "Yapılandırma hatası",
			ErrorUserMsgUnknown:       "Bilinmeyen hata oluştu",
			TaskStatusError:           "HATA",
			TaskStatusCancelled:       "İptal edildi",
			ErrFieldRequired:          "alan boş bırakılamaz",
			ErrPathEmpty:              "yol boş olamaz",
			ErrPathInvalidChar:        "yol geçersiz bir karakter içeriyor: %c",
			ErrURLEmpty:               "URL boş olamaz",
			ErrURLScheme:              "URL http:// veya https:// ile başlamalıdır",
			ErrValueEmpty:             "değer boş olamaz",
			ErrValueAlphaNumeric:      "değer yalnızca harf ve rakam içermelidir",
			ErrDefaultValueInvalid:    "varsayılan değer geçersiz",
			ErrDefaultValueEmpty:      "varsayılan değer boş",
			CancelShort:               "İptal",
			NeedSelectAtLeastOne:      "! En az bir öğe seçin",
			ScrollAboveFormat:         "%s %s %d yukarıda",
			ScrollBelowFormat:         "%s %s %d aşağıda",
			SingleSelectHelp:          "[← çıkış, ↑/↓ gezin, →/Enter seç, Q/Esc çıkış]",
			MultiSelectHelp:           "[← çıkış, ↑/↓ gezin, →/boşluk seç, Enter onay, Q/Esc çıkış]",
			MultiSelectHelpSelectAll:  "[← çıkış, ↑/↓ gezin, →/boşluk tümünü değiştir, Enter onay, Q/Esc çıkış]",
			SelectAllDefaultText:      "Tümünü seç",
			InputConfirmHint:          "[Enter onay, Ctrl+C iptal]",
			InputFormatLabel:          "Biçim:",
			InputHintPassword:         "Güçlü bir parola kullanın",
			InputHintEmail:            "Örnek: user@example.com",
			InputHintNumber:           "Bir sayı girin",
			InputHintIP:               "Örnek: 192.168.1.1",
			InputHintDomain:           "Örnek: example.com",
		},
		"be": {
			StatusSuccess:             "Паспяхова",
			StatusProblem:             "Праблема",
			StatusInProgress:          "Выконваецца",
			SummaryCompleted:          "Удала завершана",
			SummaryOf:                 "з",
			SummaryTasks:              "задач",
			DefaultNo:                 "Не",
			DefaultYes:                "Так",
			DefaultSuccessLabel:       "Гатова",
			DefaultFromSummaryLabel:   "з",
			DefaultTasksSummaryLabel:  "задач",
			DefaultErrorLabel:         "Памылка",
			DefaultCancelLabel:        "Адменена карыстальнікам",
			DefaultSelectedLabel:      "карыстальнік выбраў",
			DefaultYesLabel:           "Паспяхова",
			DefaultNoLabel:            "АДМОВА",
			TaskCancelledByUser:       "[адменена карыстальнікам]",
			TaskExitHint:              "[Для выхаду націсніце Ctrl+C]",
			DefaultPrompt:             "Увядзіце значэнне",
			PasswordMask:              '*',
			DefaultPlaceholder:        "...",
			DefaultSeparator:          ", ",
			ErrorTypeValidation:       "ВАЛІДАЦЫЯ",
			ErrorTypeUserCancel:       "АДМЕНА",
			ErrorTypeTimeout:          "ТАЙМАЎТ",
			ErrorTypeNetwork:          "СЕТКА",
			ErrorTypeFileSystem:       "ФАЙЛ",
			ErrorTypePermission:       "ДАСТУП",
			ErrorTypeConfig:           "КАНФІГ",
			ErrorTypeUnknown:          "ПАМЫЛКА",
			ErrorMsgUnknown:           "невядомая памылка",
			ErrorMsgTaskPrefix:        "задача '%s': ",
			ErrorMsgCanceled:          "адменена карыстальнікам",
			ErrorMsgTimeout:           "аперацыя не завершылася за %v",
			ErrorMsgPermission:        "недастаткова прав для доступу да %s",
			ErrorUserMsgValidation:    "Праверце карэктнасць уведзеных дадзеных",
			ErrorUserMsgCancel:        "Аперацыя адменена",
			ErrorUserMsgTimeout:       "Аперацыя заняла занадта шмат часу",
			ErrorUserMsgNetwork:       "Праблема з сеткавым злучэннем",
			ErrorUserMsgFileSystem:    "Праблема доступу да файлавай сістэмы",
			ErrorUserMsgPermission:    "Недастаткова прав для выканання аперацыі",
			ErrorUserMsgConfiguration: "Памылка ў наладках",
			ErrorUserMsgUnknown:       "Адбылася невядомая памылка",
			TaskStatusError:           "ПАМЫЛКА",
			TaskStatusCancelled:       "Адменена",
			ErrFieldRequired:          "поле павінна быць запоўнена",
			ErrPathEmpty:              "шлях не можа быць пустым",
			ErrPathInvalidChar:        "шлях утрымлівае недапушчальны сімвал: %c",
			ErrURLEmpty:               "URL не можа быць пустым",
			ErrURLScheme:              "URL павінен пачынацца з http:// або https://",
			ErrValueEmpty:             "значэнне не можа быць пустым",
			ErrValueAlphaNumeric:      "значэнне павінна змяшчаць толькі літары і лічбы",
			ErrDefaultValueInvalid:    "значэнне па змаўчанні не валіднае",
			ErrDefaultValueEmpty:      "значэнне па змаўчанні пустое",
			CancelShort:               "Адменена",
			NeedSelectAtLeastOne:      "! Неабходна выбраць хаця б адзін элемент",
			ScrollAboveFormat:         "%s %s %d вышэй",
			ScrollBelowFormat:         "%s %s %d ніжэй",
			SingleSelectHelp:          "[← выхад, ↑/↓ навігацыя, →/Enter выбар, Q/Esc — выхад]",
			MultiSelectHelp:           "[← выхад, ↑/↓ навігацыя, →/прабел выбар, Enter — пацвярджэнне, Q/Esc — выхад]",
			MultiSelectHelpSelectAll:  "[← выхад, ↑/↓ навігацыя, →/прабел пераключыць усе, Enter — пацвярджэнне, Q/Esc — выхад]",
			SelectAllDefaultText:      "Выбраць усе",
			InputConfirmHint:          "[Enter — пацвердзіць, Ctrl+C — скасаваць]",
			InputFormatLabel:          "Фармат:",
			InputHintPassword:         "Выкарыстоўвайце надзейны пароль",
			InputHintEmail:            "Прыклад: user@example.com",
			InputHintNumber:           "Увядзіце лік",
			InputHintIP:               "Прыклад: 192.168.1.1",
			InputHintDomain:           "Прыклад: example.com",
		},
		"uk": {
			StatusSuccess:             "УСПІХ",
			StatusProblem:             "ПРОБЛЕМА",
			StatusInProgress:          "ВИКОНУЄТЬСЯ",
			SummaryCompleted:          "Успішно завершено",
			SummaryOf:                 "з",
			SummaryTasks:              "завдань",
			DefaultNo:                 "Ні",
			DefaultYes:                "Так",
			DefaultSuccessLabel:       "Готово",
			DefaultFromSummaryLabel:   "з",
			DefaultTasksSummaryLabel:  "завдань",
			DefaultErrorLabel:         "Помилка",
			DefaultCancelLabel:        "Скасовано користувачем",
			DefaultSelectedLabel:      "користувач обрав",
			DefaultYesLabel:           "УСПІХ",
			DefaultNoLabel:            "ВІДМОВА",
			TaskCancelledByUser:       "[скасовано користувачем]",
			TaskExitHint:              "[Для виходу натисніть Ctrl+C]",
			DefaultPrompt:             "Введіть значення",
			PasswordMask:              '*',
			DefaultPlaceholder:        "...",
			DefaultSeparator:          ", ",
			ErrorTypeValidation:       "ВАЛІДАЦІЯ",
			ErrorTypeUserCancel:       "СКАСУВАННЯ",
			ErrorTypeTimeout:          "ТАЙМ-АУТ",
			ErrorTypeNetwork:          "МЕРЕЖА",
			ErrorTypeFileSystem:       "ФАЙЛ",
			ErrorTypePermission:       "ДОСТУП",
			ErrorTypeConfig:           "КОНФІГ",
			ErrorTypeUnknown:          "ПОМИЛКА",
			ErrorMsgUnknown:           "невідома помилка",
			ErrorMsgTaskPrefix:        "завдання '%s': ",
			ErrorMsgCanceled:          "скасовано користувачем",
			ErrorMsgTimeout:           "операцію не завершено за %v",
			ErrorMsgPermission:        "недостатньо прав для доступу до %s",
			ErrorUserMsgValidation:    "Перевірте правильність введених даних",
			ErrorUserMsgCancel:        "Операцію скасовано",
			ErrorUserMsgTimeout:       "Операція триває надто довго",
			ErrorUserMsgNetwork:       "Проблема з мережевим з'єднанням",
			ErrorUserMsgFileSystem:    "Проблема доступу до файлової системи",
			ErrorUserMsgPermission:    "Недостатньо прав для виконання операції",
			ErrorUserMsgConfiguration: "Помилка в налаштуваннях",
			ErrorUserMsgUnknown:       "Сталася невідома помилка",
			TaskStatusError:           "ПОМИЛКА",
			TaskStatusCancelled:       "Скасовано",
			ErrFieldRequired:          "поле є обов'язковим",
			ErrPathEmpty:              "шлях не може бути порожнім",
			ErrPathInvalidChar:        "шлях містить неприпустимий символ: %c",
			ErrURLEmpty:               "URL не може бути порожнім",
			ErrURLScheme:              "URL має починатися з http:// або https://",
			ErrValueEmpty:             "значення не може бути порожнім",
			ErrValueAlphaNumeric:      "значення повинно містити лише літери та цифри",
			ErrDefaultValueInvalid:    "значення за замовчуванням некоректне",
			ErrDefaultValueEmpty:      "значення за замовчуванням порожнє",
			CancelShort:               "Скасовано",
			NeedSelectAtLeastOne:      "! Потрібно вибрати принаймні один елемент",
			ScrollAboveFormat:         "%s %s %d вище",
			ScrollBelowFormat:         "%s %s %d нижче",
			SingleSelectHelp:          "[← вихід, ↑/↓ навігація, →/Enter вибір, Q/Esc — вихід]",
			MultiSelectHelp:           "[← вихід, ↑/↓ навігація, →/пробіл вибір, Enter — підтвердження, Q/Esc — вихід]",
			MultiSelectHelpSelectAll:  "[← вихід, ↑/↓ навігація, →/пробіл перемкнути всі, Enter — підтвердження, Q/Esc — вихід]",
			SelectAllDefaultText:      "Вибрати всі",
			InputConfirmHint:          "[Enter — підтвердити, Ctrl+C — скасувати]",
			InputFormatLabel:          "Формат:",
			InputHintPassword:         "Використовуйте надійний пароль",
			InputHintEmail:            "Приклад: user@example.com",
			InputHintNumber:           "Введіть число",
			InputHintIP:               "Приклад: 192.168.1.1",
			InputHintDomain:           "Приклад: example.com",
		},
	}
)

// applyDictionary применяет выбранный словарь к глобальным переменным.
func applyDictionary(dict dictionary) {
	StatusSuccess = dict.StatusSuccess
	StatusProblem = dict.StatusProblem
	StatusInProgress = dict.StatusInProgress
	SummaryCompleted = dict.SummaryCompleted
	SummaryOf = dict.SummaryOf
	SummaryTasks = dict.SummaryTasks
	DefaultNo = dict.DefaultNo
	DefaultYes = dict.DefaultYes
	DefaultSuccessLabel = dict.DefaultSuccessLabel
	DefaultFromSummaryLabel = dict.DefaultFromSummaryLabel
	DefaultTasksSummaryLabel = dict.DefaultTasksSummaryLabel
	DefaultErrorLabel = dict.DefaultErrorLabel
	DefaultCancelLabel = dict.DefaultCancelLabel
	DefaultSelectedLabel = dict.DefaultSelectedLabel
	DefaultYesLabel = dict.DefaultYesLabel
	DefaultNoLabel = dict.DefaultNoLabel
	TaskCancelledByUser = dict.TaskCancelledByUser
	TaskExitHint = dict.TaskExitHint
	DefaultPrompt = dict.DefaultPrompt
	PasswordMask = dict.PasswordMask
	DefaultPlaceholder = dict.DefaultPlaceholder
	DefaultSeparator = dict.DefaultSeparator
	ErrorTypeValidation = dict.ErrorTypeValidation
	ErrorTypeUserCancel = dict.ErrorTypeUserCancel
	ErrorTypeTimeout = dict.ErrorTypeTimeout
	ErrorTypeNetwork = dict.ErrorTypeNetwork
	ErrorTypeFileSystem = dict.ErrorTypeFileSystem
	ErrorTypePermission = dict.ErrorTypePermission
	ErrorTypeConfig = dict.ErrorTypeConfig
	ErrorTypeUnknown = dict.ErrorTypeUnknown
	ErrorMsgUnknown = dict.ErrorMsgUnknown
	ErrorMsgTaskPrefix = dict.ErrorMsgTaskPrefix
	ErrorMsgCanceled = dict.ErrorMsgCanceled
	ErrorMsgTimeout = dict.ErrorMsgTimeout
	ErrorMsgPermission = dict.ErrorMsgPermission
	ErrorUserMsgValidation = dict.ErrorUserMsgValidation
	ErrorUserMsgCancel = dict.ErrorUserMsgCancel
	ErrorUserMsgTimeout = dict.ErrorUserMsgTimeout
	ErrorUserMsgNetwork = dict.ErrorUserMsgNetwork
	ErrorUserMsgFileSystem = dict.ErrorUserMsgFileSystem
	ErrorUserMsgPermission = dict.ErrorUserMsgPermission
	ErrorUserMsgConfiguration = dict.ErrorUserMsgConfiguration
	ErrorUserMsgUnknown = dict.ErrorUserMsgUnknown
	TaskStatusError = dict.TaskStatusError
	TaskStatusCancelled = dict.TaskStatusCancelled
	ErrFieldRequired = dict.ErrFieldRequired
	ErrPathEmpty = dict.ErrPathEmpty
	ErrPathInvalidChar = dict.ErrPathInvalidChar
	ErrURLEmpty = dict.ErrURLEmpty
	ErrURLScheme = dict.ErrURLScheme
	ErrValueEmpty = dict.ErrValueEmpty
	ErrValueAlphaNumeric = dict.ErrValueAlphaNumeric
	ErrDefaultValueInvalid = dict.ErrDefaultValueInvalid
	ErrDefaultValueEmpty = dict.ErrDefaultValueEmpty
	CancelShort = dict.CancelShort
	NeedSelectAtLeastOne = dict.NeedSelectAtLeastOne
	ScrollAboveFormat = dict.ScrollAboveFormat
	ScrollBelowFormat = dict.ScrollBelowFormat
	SingleSelectHelp = dict.SingleSelectHelp
	MultiSelectHelp = dict.MultiSelectHelp
	MultiSelectHelpSelectAll = dict.MultiSelectHelpSelectAll
	SelectAllDefaultText = dict.SelectAllDefaultText
	InputConfirmHint = dict.InputConfirmHint
	InputFormatLabel = dict.InputFormatLabel
	InputHintPassword = dict.InputHintPassword
	InputHintEmail = dict.InputHintEmail
	InputHintNumber = dict.InputHintNumber
	InputHintIP = dict.InputHintIP
	InputHintDomain = dict.InputHintDomain
}

// SetLanguage обновляет текущий язык и возвращает фактически установленное значение.
func SetLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	if dict, ok := dictionaries[lang]; ok {
		applyDictionary(dict)
		currentLanguage = lang
		return lang
	}
	applyDictionary(dictionaries["en"])
	currentLanguage = "en"
	return currentLanguage
}

// CurrentLanguage возвращает код текущего языка.
func CurrentLanguage() string {
	return currentLanguage
}

// SupportedLanguages возвращает список поддерживаемых языков.
func SupportedLanguages() []string {
	keys := make([]string, 0, len(dictionaries))
	for k := range dictionaries {
		keys = append(keys, k)
	}
	return keys
}

// SetDefaultLanguage задаёт язык по умолчанию и возвращает фактически применённое значение.
func SetDefaultLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	dict, ok := dictionaries[lang]
	if !ok {
		dict = dictionaries["en"]
		lang = "en"
	}
	defaultLanguage = lang
	applyDictionary(dict)
	currentLanguage = lang
	return lang
}

func init() {
	if dict, ok := dictionaries[defaultLanguage]; ok {
		applyDictionary(dict)
		currentLanguage = defaultLanguage
	} else {
		applyDictionary(dictionaries["en"])
		currentLanguage = "en"
	}
}
