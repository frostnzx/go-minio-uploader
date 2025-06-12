"use client";
import { Button, DatePicker, Form, FormProps, Input } from "antd";
import React, { ChangeEvent, useState } from "react";
import type { UploadProps } from "antd";
import { message, Upload } from "antd";
import { InboxOutlined } from "@ant-design/icons";
import dayjs from "dayjs";

const { Dragger } = Upload;

type FieldType = {
    name: string;
    description: string;
    date: dayjs.Dayjs;
};

export default function InputForm({
    onSubmit,
    setFormData,
    formData,
    fileList,
    setFileList,
}: {
    onSubmit: Function;
    setFormData: Function;
    formData: any;
    fileList: File[];
    setFileList: Function;
}) {
    const props: UploadProps = {
        beforeUpload: (file) => {
            setFileList((prev: any) => [...prev, file]);
            return false; // Prevent auto upload
        },
        onRemove: (file) => {
            setFileList((prev: any) =>
                prev.filter((f: any) => f.uid !== file.uid)
            );
        },
        fileList: fileList as any, // required to control Upload component
    };

    // Form functions
    const onInputChange = (e: ChangeEvent<HTMLInputElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };
    // For DatePicker
    const onDateChange = (date: dayjs.Dayjs | null) => {
        setFormData({ ...formData, date });
    };

    const onFinish: FormProps<FieldType>["onFinish"] = (values) => {
        onSubmit();
    };
    const onFinishFailed: FormProps<FieldType>["onFinishFailed"] = (
        errorInfo
    ) => {
        console.log("Failed:", errorInfo);
    };

    return (
        <>
            <Form
                name="basic"
                labelCol={{ span: 8 }}
                wrapperCol={{ span: 16 }}
                style={{ maxWidth: 600, justifyContent: "start" }}
                initialValues={{ remember: true }}
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
            >
                <Form.Item<FieldType>
                    label="Name"
                    name="name"
                    rules={[
                        {
                            required: true,
                            message: "Please input your username!",
                        },
                    ]}
                >
                    <Input name="name" onChange={onInputChange} />
                </Form.Item>
                <Form.Item<FieldType>
                    label="Description"
                    name="description"
                    rules={[
                        {
                            required: true,
                            message: "Please input your description!",
                        },
                        {
                            max: 25,
                            message:
                                "Description must be less than 25 characters",
                        },
                    ]}
                >
                    <Input name="description" onChange={onInputChange} />
                </Form.Item>
                <Form.Item<FieldType>
                    label="Date and time : "
                    name="date"
                    rules={[
                        {
                            type: "object",
                            required: true,
                            message: "Please input the date and time",
                        },
                    ]}
                >
                    <DatePicker
                        showTime
                        name="date"
                        onChange={onDateChange}
                        format="DD/MM/YYYY HH:mm"
                        placeholder="DD/MM/YYYY HH:mm"
                        allowClear={false}
                    />
                </Form.Item>
                <Dragger {...props}>
                    <p className="ant-upload-drag-icon">
                        <InboxOutlined />
                    </p>
                    <p className="ant-upload-text">
                        Click or drag file to this area to upload
                    </p>
                    <p className="ant-upload-hint">
                        Support for a single or bulk upload. Strictly prohibited
                        from uploading company data or other banned files.
                    </p>
                </Dragger>

                <Form.Item
                    label={null}
                    wrapperCol={{ span: 24 }}
                    style={{
                        marginTop: "50px",
                        textAlign: "center",
                    }}
                >
                    <Button type="primary" htmlType="submit">
                        Submit
                    </Button>
                </Form.Item>
            </Form>
        </>
    );
}
