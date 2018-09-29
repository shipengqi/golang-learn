module.exports = {
    base: '/go-learn-demo/',
    title: 'Go 语言学习笔记',
    description: 'Go 数据类型，函数，方法，接口，并发编程，高级编程',
    head: [
        ['link', {
            rel: 'icon',
            href: 'https://upload.wikimedia.org/wikipedia/commons/2/23/Go_Logo_Aqua.svg'
        }]
    ],
    markdown: {
        toc: {
            includeLevel: [2, 3, 4, 5, 6]
        }
    },
    themeConfig: {
        repo: 'shipengqi/go-learn-demo',
        docsDir: 'docs',
        editLinks: true,
        editLinkText: '错别字纠正',
        sidebarDepth: 3,
        nav: [{
            text: 'Go 语言入门',
            link: '/go_basic/',
        }, {
            text: '并发编程',
            link: '/concurrent/'
        }, {
            text: '高级编程',
            link: '/advanced/'
        }],
        sidebar: {
            '/go_basic/': [{
                title: 'Go 语言入门（持续更新中...）',
                children: [
                    '',
                    'basic_data_types',
                    'basic_syntax',
                    'function',
                    'OOP',
                    'test',
                    'reflect'
                ]
            }],
            '/concurrent/': [{
                title: '并发编程',
                children: [
                    ''
                ]
            }],
            '/advanced/': [{
                title: '高级编程',
                children: [
                    ''
                ]
            }]
        }
    }
}