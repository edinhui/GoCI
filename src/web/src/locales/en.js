export default {
  sidebar: {
    schemaManagement: 'Schema Management',
    configEditor: 'Config Editor'
  },
  schemaEditor: {
    title: 'JSON Schema Editor',
    editTitle: 'Edit Schema: {name}',
    properties: 'Schema Properties',
    preview: 'Schema Preview',
    addProperty: 'Add Property',
    exportSchema: 'Export Schema',
    save: 'Save Schema',
    backToList: 'Back to List',
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
  },
  schemaViewer: {
    title: 'Schema Viewer',
    schemaDetails: 'Schema Details',
    backToList: 'Back to List',
    edit: 'Edit',
    id: 'ID',
    name: 'Name',
    description: 'Description',
    createdAt: 'Created At',
    updatedAt: 'Updated At',
    schemaContent: 'Schema Content',
    configExample: 'Config Example',
    loadError: 'Failed to load schema'
  },
  schemaList: {
    title: 'Schema List',
    availableSchemas: 'Available Schemas',
    search: 'Search Schema',
    createNew: 'Create New Schema',
    id: 'ID',
    name: 'Name',
    description: 'Description',
    updatedAt: 'Updated At',
    actions: 'Actions',
    view: 'View',
    edit: 'Edit',
    delete: 'Delete',
    deleteConfirm: 'Confirm delete this Schema?',
    noSchemas: 'No Schemas',
    loadError: 'Failed to load schema list'
  }
}
