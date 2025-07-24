export default {
  schemaEditor: {
    title: 'JSON Schema Editor',
    properties: 'Schema Properties',
    preview: 'Schema Preview',
    addProperty: 'Add Property',
    exportSchema: 'Export Schema',
    validationRules: 'Validation Rules'
  },
  propertyDialog: {
    add: 'Add Property',
    edit: 'Edit Property',
    name: 'Property Name',
    type: 'Property Type',
    required: 'Required',
    description: 'Description',
    fixedFieldValue: 'Fixed Field Value',
    fixedFieldHint: 'This value will be used as default and cannot be modified in configuration',
    cancel: 'Cancel',
    save: 'Save'
  },
  propertyTypes: {
    string: 'String',
    number: 'Number',
    boolean: 'Boolean',
    object: 'Object',
    array: 'Array'
  },
  validation: {
    string: {
      minLength: 'Min Length',
      maxLength: 'Max Length',
      pattern: 'Pattern',
      patternHint: 'Example: ^[a-zA-Z0-9]+$',
      format: 'Format'
    },
    number: {
      minimum: 'Minimum',
      maximum: 'Maximum',
      multipleOf: 'Multiple Of',
      exclusiveMinimum: 'Exclusive Min',
      exclusiveMaximum: 'Exclusive Max'
    },
    array: {
      minItems: 'Min Items',
      maxItems: 'Max Items',
      uniqueItems: 'Unique Items'
    },
    enum: {
      title: 'Enum Values',
      addValue: '+ Add Value'
    },
    formats: {
      email: 'Email',
      uri: 'URI',
      date: 'Date',
      dateTime: 'Date-Time',
      hostname: 'Hostname',
      ipv4: 'IPv4',
      ipv6: 'IPv6'
    }
  },
  fixedFields: {
    title: 'Fixed Field',
    readOnly: 'ReadOnly'
  },
  language: {
    switch: 'Switch to 中文',
    current: 'English'
  }
}
