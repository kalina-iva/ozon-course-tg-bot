package messages

const manual = `Привет! Это дневник расходов.
Описание команд:
/newcat {name} - добавление новой категории
/allcat - просмотр всех категории
/newexpense {categoryNumber} {amount} {date} - добавление нового расхода. Если дата не указана, используется текущая
/report {y|m|w} - получение отчета за последний год/месяц/неделю
`

const (
	unknownCommand        = "Неизвестная команда"
	noCategories          = "Нет категорий"
	expenseAdded          = "Расход добавлен"
	needCategoryAndAmount = "Необходимо указать категорию и сумму"
	invalidCategoryNumber = "Некорректный номер категории"
	invalidAmount         = "Некорректная сумма расхода"
	invalidDate           = "Некорректная дата"
	needPeriod            = "Необходимо указать период"
	invalidPeriod         = "Некорректный период"
	categoryNotFound      = "Не найдена категория с номером %s"
	categoryCreated       = "Создана категория %s. Ее номер %v"
	noCategoryName        = "Нет названия категории"
)
