#import "conf.typ": conf, intro, conclusion
#show: conf.with(
  title: [Отчет по теории графов],
  type: "pract",
  info: (
    author: (
      name: [Толстова Роберта Сергеевича],
      faculty: [КНиИТ],
      group: "351",
      sex: "male",
    ),
    inspector: (
      degree: "доцент, к. ф.-м. н.",
      name: "С. В. Миронов",
    ),
  ),
  settings: (
    title_page: (
      enabled: true,
    ),
    contents_page: (
      enabled: true,
    ),
  ),
)

//#intro
//#conclusion
#for value in ("1", "2", "3", "4") {
  include "sections/" + value + ".typ"
}