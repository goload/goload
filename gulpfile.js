'use strict';

var gulp = require('gulp');  // Base gulp package
var babelify = require('babelify'); // Used to convert ES6 & JSX to ES5
var browserify = require('browserify'); // Providers "require" support, CommonJS
var notify = require('gulp-notify'); // Provides notification to both the console and Growel
var rename = require('gulp-rename'); // Rename sources
var sourcemaps = require('gulp-sourcemaps'); // Provide external sourcemap files
var livereload = require('gulp-livereload'); // Livereload support for the browser
var gutil = require('gulp-util'); // Provides gulp utilities, including logging and beep
var chalk = require('chalk'); // Allows for coloring for logging
var source = require('vinyl-source-stream'); // Vinyl stream support
var buffer = require('vinyl-buffer'); // Vinyl stream support
var merge = require('merge-stream'); // Object merge tool
var duration = require('gulp-duration'); // Time aspects of your gulp process
var concat = require('gulp-concat');
var browserifyInc = require('browserify-incremental')
var cleanCSS = require('gulp-clean-css');
var del = require('del');
var uglify = require('gulp-uglify');
var run = require('gulp-run');
var rename = require('gulp-rename');
var zip = require('gulp-zip');
var version= require('./package.json').version;

// Configuration for Gulp
var config = {
    js: {
        src: './client/js/app.jsx',
        watch: './client/js/**/*',
        outputDir: './public/',
        outputFile: 'app.js',
    },
};

// Error reporting function
function mapError(err) {
    if (err.fileName) {
        // Regular error
        gutil.log(chalk.red(err.name)
            + ': ' + chalk.yellow(err.fileName.replace(__dirname + '/src/js/', ''))
            + ': ' + 'Line ' + chalk.magenta(err.lineNumber)
            + ' & ' + 'Column ' + chalk.magenta(err.columnNumber || err.column)
            + ': ' + chalk.blue(err.description));
    } else {
        // Browserify error..
        gutil.log(chalk.red(err.name)
            + ': '
            + chalk.yellow(err.message));
    }
}

// Completes the final file outputs
gulp.task('js', function () {
    var bundler = browserify(config.js.src, browserifyInc.args);
    bundler.transform(babelify, {presets: ['es2015', 'react']});
    browserifyInc(bundler, {cacheFile: './browserify-cache.json'})
    var bundleTimer = duration('Javascript bundle time');

    bundler
        .bundle()
        .on('error', mapError) // Map error reporting
        .pipe(source('main.jsx')) // Set source name
        .pipe(buffer()) // Convert to gulp pipeline
        .pipe(rename(config.js.outputFile)) // Rename the output file
        .pipe(sourcemaps.init({loadMaps: true})) // Extract the inline sourcemaps
        .pipe(sourcemaps.write('./map')) // Set folder for sourcemaps to output to
        .pipe(gulp.dest(config.js.outputDir)) // Set the output folder
        .pipe(notify({
            message: 'Generated file: <%= file.relative %>',
        })) // Output the file being created
        .pipe(bundleTimer) // Output time timing of the file creation
        .pipe(livereload()); // Reload the view in the browser
});

gulp.task('css', function () {
    return gulp.src(['client/css/**/*.css','node_modules/bootstrap/dist/css/bootstrap.min.css'])
        .pipe(concat('style.css'))
        .pipe(gulp.dest(config.js.outputDir))
});

gulp.task('buildjs', function () {
    process.env.NODE_ENV = 'production';
    var bundler = browserify(config.js.src);
    bundler.transform(babelify, {presets: ['es2015', 'react']});
    var bundleTimer = duration('Javascript bundle time');

    return bundler
        .bundle()
        .on('error', mapError) // Map error reporting
        .pipe(source('main.jsx')) // Set source name
        .pipe(buffer()) // Convert to gulp pipeline
        .pipe(rename(config.js.outputFile)) // Rename the output file
        .pipe(uglify())
        .pipe(gulp.dest('dist/public/')) // Set the output folder
        .pipe(bundleTimer); // Output time timing of the file creation
});

gulp.task('buildcss', function () {
    return gulp.src(['node_modules/bootstrap/dist/css/bootstrap.min.css', 'client/css/**/*.css'])
        .pipe(concat('style.css'))
        .pipe(cleanCSS())
        .pipe(gulp.dest('dist/public/'))
});

gulp.task('watch', function () {
    gulp.watch('client/css/**/*.css').on('change', gulp.series('css'));
    gulp.watch('client/js/**/*.jsx').on('change', gulp.series('js'));
});

// Gulp task for build
gulp.task('default', gulp.parallel('watch', 'js', 'css', function () {
    livereload.listen(); // Start livereload server
}));

gulp.task('cleanDist', function (done) {
    del.sync(['./dist']);
    done();
});

gulp.task('buildAllGo', gulp.series('cleanDist', function (done) {
    return merge(
        run('env GOOS=linux GOARCH=arm go build -ldflags "-X main.Version='+version+'" -o dist/linux-arm/goload goload/server').exec(),
        run('env GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version='+version+'" -o dist/linux-arm64/goload goload/server').exec(),
        run('env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version='+version+'" -o dist/linux-amd64/goload goload/server').exec(),
        run('env GOOS=linux GOARCH=386 go build -ldflags "-X main.Version='+version+'" -o dist/linux-386/goload goload/server').exec(),
        //windows
        run('env GOOS=windows GOARCH=386 go build -ldflags "-X main.Version='+version+'" -o dist/windows-386/goload goload/server').exec(),
        run('env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version='+version+'" -o dist/windows-amd64/goload goload/server').exec(),
        //osx
        run('env GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version='+version+'" -o dist/osx-amd64/goload goload/server').exec(),
        run('env GOOS=darwin GOARCH=386 go build -ldflags "-X main.Version='+version+'" -o dist/osx-386/goload goload/server').exec()
    );
    //env GOOS=darwin GOARCH=arm64 go build -o build/osx-arm64/server goload/server
    //env GOOS=darwin GOARCH=arm go build -o build/osx-arm/server goload/server
}));

gulp.task('copyPublic', function () {
    return merge(
        gulp.src('./public/fonts/**/*').pipe(gulp.dest('./dist/public/fonts')),
        gulp.src('./public/images/**/*').pipe(gulp.dest('./dist/public/images')),
        gulp.src('./public/index.html').pipe(gulp.dest('./dist/public')),
        gulp.src('./defaultConfig.json').pipe(gulp.dest('./dist'))
    );
});

gulp.task('build', gulp.series('cleanDist',gulp.parallel('buildAllGo', 'buildjs', 'buildcss', 'copyPublic'), function (done) {
    del.sync(['goload*.zip']);
    gulp.src('dist/**/*')
        .pipe(zip('goload-v'+version+'.zip'))
        .pipe(gulp.dest('.'));
    done();
}));
