extern crate env_logger;
#[macro_use]
extern crate log;


fn init_logger() {
    let mut builder = env_logger::Builder::new();
    builder.format(|buf, record| {
        use std::io::Write;
        writeln!(buf, "{} {} {}:{} {}", 
                 buf.timestamp(), 
                 record.level(), 
                 record.file().unwrap_or("?"), 
                 record.line().unwrap_or(0), 
                 record.args())                
    });
    builder.parse("info");
    builder.init();
}

fn do_something() {
    info!("Generating types...");
    trace!("Working on {} ...", definition_path);
}
