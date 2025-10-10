package processor

import (
	"arithmetic-calc/internal/expression"
	"encoding/json"
)

func ProcessJSON(data []byte) ([]byte, error) {
	// interface{} -пустой интерфейс, обозначает любой тип, как any в typescript
	// данные в json могут быть любого типа.
	var parsedJson interface{}

	// выполняем код и сразу проверяем ошибку
	// Это называется if с коротким обьявлением / short statement
	// Обьеденияет 2 действия: иницилизацию и проверку условия
	// Главная причина в области видимости
	// переменная err видна только внутри блока if {}
	// У if есть 2 режима:
	// if <условие> {}
	// if <инициализация>; <условие> {} - if с коротким обьявлением

	if err := json.Unmarshal(data, &parsedJson); err != nil {
		return nil, err
	}

	// запускаем рекурсию на корневом узле карты json
	processedJson := recursiveProcessString(parsedJson)

	return json.MarshalIndent(processedJson, "", " ")
}

// функция принимает любой тип
// в нашем случае это корень json в виде map
// проверка на map запускает цикл и рекурсивно проходит по строкам
// если строка значение, а не ключ то мы применяем функцию recursionProcessString

func recursiveProcessString(node interface{}) interface{} {
	// этот синтаксис работает только в switch
	// это не обычное присваивание
	// v - не хранит тип (в go нельзя хранить тип)
	// Это называется type switch и работает только в switch
	// .(type) - ключевое слово для работы с interface{}
	// оно работает только в switch. Заглянет и покажет тип.
	// .(type) - это ключевое слово работает только с интерфейсами
	// необязательно для type именно interface{} может быть любой интерфейс
	// и чтобы определить к какому типу относится этот интерфейс из реализаций

	// v := - это не одна и таже переменная.
	// Для каждого блока case компилятор неявно создает переменную с именем v
	// и преобразует безопасно к этому типу.

	// То есть V для каждого кейс своя переменная типа case.

	switch v := node.(type) {
	case map[string]interface{}:
		// v - будет типа map[string]interface{}
		// и компилятор преобразует interface{} в этот тип.
		// проходимся по map и для каждого значения
		// и обрабатываем знание, не ключ
		for key, value := range v {
			v[key] = recursiveProcessString(value)
		}
		return v
	case []interface{}:
		// v - будет типа []interface{}
		// проходимся по каждому элементу среза
		// и обрабатываем его
		for i, item := range v {
			v[i] = recursiveProcessString(item)
		}
		// если все обработали возвращаем сам срез
		return v
	case string:
		processedString, err := expression.ProcessString(v)
		if err != nil {
			return v
		}
		return processedString
	default:
		return v
	}
}
