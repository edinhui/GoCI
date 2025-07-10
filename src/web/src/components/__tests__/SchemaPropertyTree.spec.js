import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SchemaPropertyTree from '../SchemaPropertyTree.vue'

// Mock the Element Plus components
vi.mock('element-plus', () => ({
  ElTree: {
    name: 'ElTree',
    template: '<div class="el-tree"><slot :node="{}" :data="data[0]"></slot></div>',
    props: ['data']
  },
  ElTag: {
    name: 'ElTag',
    template: '<div class="el-tag"><slot></slot></div>',
    props: ['type', 'size']
  },
  ElButton: {
    name: 'ElButton',
    template: '<button class="el-button"><slot></slot></button>',
    props: ['size', 'type']
  },
  ElButtonGroup: {
    name: 'ElButtonGroup',
    template: '<div class="el-button-group"><slot></slot></div>'
  },
  ElDialog: {
    name: 'ElDialog',
    template: '<div class="el-dialog"><div class="header"><slot name="header">{{title}}</slot></div><div class="body"><slot></slot></div><div class="footer"><slot name="footer"></slot></div></div>',
    props: ['modelValue', 'title']
  },
  ElForm: {
    name: 'ElForm',
    template: '<form class="el-form"><slot></slot></form>',
    props: ['model', 'labelWidth']
  },
  ElFormItem: {
    name: 'ElFormItem',
    template: '<div class="el-form-item"><label v-if="label">{{label}}</label><slot></slot></div>',
    props: ['label']
  },
  ElInput: {
    name: 'ElInput',
    template: '<input class="el-input" />',
    props: ['modelValue']
  },
  ElSelect: {
    name: 'ElSelect',
    template: '<select class="el-select"><slot></slot></select>',
    props: ['modelValue', 'placeholder']
  },
  ElOption: {
    name: 'ElOption',
    template: '<option class="el-option" :value="value">{{label}}</option>',
    props: ['label', 'value']
  },
  ElSwitch: {
    name: 'ElSwitch',
    template: '<input type="checkbox" class="el-switch" />',
    props: ['modelValue']
  }
}))

// Mock the Element Plus icons
vi.mock('@element-plus/icons-vue', () => ({
  Edit: { template: '<div class="el-icon-edit"></div>' },
  Plus: { template: '<div class="el-icon-plus"></div>' },
  Delete: { template: '<div class="el-icon-delete"></div>' }
}))

describe('SchemaPropertyTree', () => {
  const createWrapper = (props = {}) => {
    return mount(SchemaPropertyTree, {
      props: {
        properties: [
          {
            id: 'prop_1',
            name: 'testProperty',
            type: 'string',
            required: true,
            children: []
          }
        ],
        ...props
      },
      global: {
        stubs: {
          ElIcon: true
        }
      }
    })
  }

  it('renders property name and type correctly', () => {
    const wrapper = createWrapper()
    expect(wrapper.find('.property-name').text()).toBe('testProperty')
    expect(wrapper.find('.el-tag').text()).toContain('string')
  })

  it('emits update:property event when editing a property', async () => {
    const wrapper = createWrapper()
    
    // Click edit button
    await wrapper.find('.el-button').trigger('click')
    
    // Find save button in dialog and click it
    const saveButton = wrapper.find('.el-dialog .footer button:last-child')
    await saveButton.trigger('click')
    
    // Check if update:property event was emitted
    expect(wrapper.emitted('update:property')).toBeTruthy()
  })

  it('emits delete:property event when deleting a property', async () => {
    const wrapper = createWrapper()
    
    // Find delete button (third button) and click it
    const buttons = wrapper.findAll('.el-button')
    await buttons[2].trigger('click')
    
    // Check if delete:property event was emitted with correct ID
    expect(wrapper.emitted('delete:property')).toBeTruthy()
    expect(wrapper.emitted('delete:property')[0]).toEqual(['prop_1'])
  })
})
