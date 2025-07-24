export default {
  schemaEditor: {
    title: 'JSON Schema 编辑器',
    properties: 'Schema 属性',
    preview: 'Schema 预览',
    addProperty: '添加属性',
    exportSchema: '导出 Schema',
    validationRules: '验证规则'
  },
  propertyDialog: {
    add: '添加属性',
    edit: '编辑属性',
    name: '属性名称',
    type: '属性类型',
    required: '是否必填',
    description: '描述',
    fixedFieldValue: '固定字段值',
    fixedFieldHint: '此值将作为默认值且在配置中不可修改',
    cancel: '取消',
    save: '保存'
  },
  propertyTypes: {
    string: '字符串',
    number: '数字',
    boolean: '布尔值',
    object: '对象',
    array: '数组'
  },
  validation: {
    string: {
      minLength: '最小长度',
      maxLength: '最大长度',
      pattern: '模式',
      patternHint: '示例：^[a-zA-Z0-9]+$',
      format: '格式'
    },
    number: {
      minimum: '最小值',
      maximum: '最大值',
      multipleOf: '倍数',
      exclusiveMinimum: '严格最小值',
      exclusiveMaximum: '严格最大值'
    },
    array: {
      minItems: '最小项数',
      maxItems: '最大项数',
      uniqueItems: '唯一项'
    },
    enum: {
      title: '枚举值',
      addValue: '+ 添加枚举值'
    },
    formats: {
      email: '邮箱',
      uri: 'URI',
      date: '日期',
      dateTime: '日期时间',
      hostname: '主机名',
      ipv4: 'IPv4',
      ipv6: 'IPv6'
    }
  },
  fixedFields: {
    title: '固定字段',
    readOnly: '只读'
  },
  language: {
    switch: 'Switch to English',
    current: '中文'
  }
}
