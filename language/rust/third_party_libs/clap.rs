#[macro_use]
extern crate clap; 

use clap::{App, Arg};

fn main() {
    let matches = App::new("snamed")
        .arg(Arg::with_name("server")
             .help("server address")
             .short("s")
             .long("server")
             .required(true)
             .takes_value(true))
        .arg(Arg::with_name("next")
             .help("next server address")
             .short("n")
             .long("next")
             .takes_value(true))
        .arg(Arg::with_name("delay")
             .help("random delay")
             .short("d")
             .long("delay"))
        .arg(Arg::with_name("level")
             .help("name hierachy level")
             .short("l")
             .long("level")
             .required(true)
             .takes_value(true))
        .get_matches();
    let mut addr = matches.value_of("server").unwrap().to_string();
    addr.push_str(":53");
    let addr = addr.parse::<SocketAddr>().unwrap();

    let level = matches.value_of("level").unwrap().parse::<u8>().unwrap();
    let random_delay = matches.is_present("delay"),
}
