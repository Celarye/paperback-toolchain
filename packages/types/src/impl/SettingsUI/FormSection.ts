type SectionInfo = {
  id: string;
  header?: string;
  footer?: string;
};

export function Section(
  params: string | SectionInfo,
  items: (Application.FormItemElement<unknown> | undefined)[],
): Application.FormSectionElement {
  let info: SectionInfo;
  if (typeof params === "string") {
    info = { id: params };
  } else {
    info = params;
  }

  return {
    ...info,
    items: items.filter(
      (x) => x,
    ) as Application.FormItemElement<unknown>[],
  };
}

type ListSectionProps = {
  items: unknown;
  allowDeletion: boolean;
  onRemove: SelectorID<() => Promise<void>>;
  allowAddition: boolean;
  onAdd: SelectorID<() => Promise<void>>;
  rowBuilder: (item: unknown) => Application.FormItemElement<unknown>;
};

function ListSection(id: string, props: ListSectionProps) {
  // TODO
  // ListSection('mySection', {
  //     items: [{ value: 'hello', id: 'world' }],
  //     allowDeletion: true,
  //     onRemove: Application.selector(this, 'myItemDidRemove'),
  //     allowAddition: true,
  //     onAdd: Application.selector(this, 'myItemDidAdd'),
  //     rowBuilder: (element) => InputRow('myRow', {
  //         id: element.id,
  //         value: element.value,
  //         placeholder: 'Foo'
  //     })
  // })
}

declare global {
  namespace Application {
    interface FormSectionElement {
      id: string;
      header?: string;
      footer?: string;
      items: FormItemElement<unknown>[];
    }
  }
}
