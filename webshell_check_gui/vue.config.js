module.exports = {
    pluginOptions: {
        electronBuilder: {
            nodeIntegration: true,
            builderOptions:{
                extraFiles:[
                    {
                        from: 'gui_go/webshell_gui_go',
                        to: 'Resources/bin/webshell_gui_go',
                        filter: ['**/*']
                    }
                ],
            }
        }
    }
}