import path from 'path';
import {LOG_FOLDER} from './constants';

const {createLogger, transports, format} = require('winston');
require('winston-daily-rotate-file');

const level = process.env.NODE_ENV === 'development' ? 'debug' : 'info';

type WinstonObject = {
    timestamp: string;
    level: string;
    message: string;
}

const logger = createLogger({
    level,
    format: format.combine(
        format.colorize({all: true}),
        format.timestamp({format: 'YYYY-DD-MM HH:mm:ss'}),
        format.printf((data: WinstonObject) => `[${data.timestamp}] ${data.level}: ${data.message}`)
    ),
    transports: [
        new transports.DailyRotateFile({
            name: 'file#info',
            level: 'info',
            filename: path.join(LOG_FOLDER, '%DATE%.info.log'),
            datePattern: 'yyyy-MM-DD',
            prepend: true,
        }),
        new transports.DailyRotateFile({
            name: 'file#error',
            level: 'error',
            filename: path.join(LOG_FOLDER, '%DATE%.error.log'),
            datePattern: 'yyyy-MM-DD',
            prepend: true,
            handleExceptions: true,
        }),
        new transports.Console({
            level,
            // format: format.combine(
            //     format.timestamp(),
            //     format.colorize(),
            //     format.simple(),
            // ),
        }),
    ],
});

export default logger;

/*
// Automatically remove old log files

var fs = require('fs');
var path = require("path");
var CronJob = require('cron').CronJob;
var _ = require("lodash");
var logger = require("./logger");

var job = new CronJob('00 00 00 * *', function(){
    // Runs every day
    // at 00:00:00 AM.
    fs.readdir(path.join("/var", "logs", "somepath"), function(err, files){
        if(err){
            logger.error("error reading log files");
        } else{
            var currentTime = new Date();
            var weekFromNow = currentTime -
                (new Date().getTime() - (7 * 24 * 60 * 60 * 1000));
            _(files).forEach(function(file){
                var fileDate = file.split(".")[2]; // get the date from the file name
                if(fileDate){
                    fileDate = fileDate.replace(/-/g,"/");
                    var fileTime = new Date(fileDate);
                    if((currentTime - fileTime) > weekFromNow){
                        console.log("delete file",file);
                        fs.unlink(path.join("/var", "log", "ironbeast", file),
                            function (err) {
                                if (err) {
                                    logger.error(err);
                                }
                                logger.info("deleted log file: " + file);
                            });
                    }
                }
            });
        }
        */
