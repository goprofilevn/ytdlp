import React, {useEffect, useState} from 'react'
import {Button, Form, Input, message, TimePicker} from 'antd'
import { SetupResources, StartDownload} from '../wailsjs/go/main/App'
import {EventsOn, EventsOff} from '../wailsjs/runtime'
import dayjs from 'dayjs'
import './App.scss'

function App() {
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)
  const [isReady, setIsReady] = useState(false)


  const handleStart = async () => {
    try {
      const values = form.getFieldsValue()
      setLoading(true)
      const time = values.time.map((t: dayjs.Dayjs) => t.format('HH:mm:ss'))
      await StartDownload(values.url, {
        start: time[0],
        end: time[1],
      })
      form.setFieldsValue({
        lastDownload: values.url,
      })
    } catch (ex: any) {
      message.open({
        type: 'error',
        content: ex.message,
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    SetupResources()

    // Resource
    EventsOn('resource-start', () => {
      message.open({
        type: 'info',
        content: 'Start initializing resource',
      })
    })
    EventsOn('resource-progress', (data: {
      key: string
      title: string
      description: string
      progress: number
    }) => {
      message.open({
        type: 'info',
        key: data.key,
        content: `${data.title}: ${data.description} (${data.progress}%)`,
        duration: 0
      })
    })
    EventsOn('resource-error', (err: any) => {
      message.destroy('resource-progress')
      message.open({
        type: 'error',
        content: err.message,
      })
    })

    EventsOn('resource-finish', () => {
      setIsReady(true)
      message.destroy()
      message.open({
        type: 'success',
        content: 'Resource is ready',
      })
    })

    // Download
    EventsOn('download-start', () => {
      message.open({
        type: 'info',
        content: 'Start downloading',
      })
    })

    return () => {
      EventsOff('resource-finish')
      EventsOff('resource-error')
      EventsOff('resource-progress')
      EventsOff('resource-start')
    }
  }, [])

  return (
    <div id="App">
      <Form
        layout="vertical"
        form={form}
      >
        <Form.Item
          label="Url"
          name="url"
          rules={[{required: true, message: 'Missing url'}]}
        >
          <Input placeholder="https://www.twitch.tv/videos/2229290131"/>
        </Form.Item>
        <Form.Item
          label="Time"
          name="time"
          rules={[
            {
              required: true,
              message: "Please input time",
            },
          ]}
        >
          <TimePicker.RangePicker needConfirm={false} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" loading={loading} disabled={!isReady} onClick={handleStart}>
            Start
          </Button>
        </Form.Item>
        <Form.Item
          label="Last Download"
          name="lastDownload"
        >
          <Input disabled />
        </Form.Item>
      </Form>
    </div>
  )
}

export default App
