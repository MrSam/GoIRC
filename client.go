package main

import "net"
import "fmt"
import "bufio"
import "strings"

var conn, _ = net.Dial("tcp", "91.212.186.20:6667") 

// this is the preffered nickname, check what we get back
var my_nickname = "MrRobot"

func main() {
  raw := bufio.NewReader(conn)
  sendtoserver("NICK " + my_nickname)
  sendtoserver("USER Robot robot my.hostname :I'm a Robot")
  sendtoserver("JOIN #sam")

  addcallbacks()

  for {
    rawmessage, _ := raw.ReadString('\n')
    parsefromserver(rawmessage)
  }
} 

func sendtoserver(message string) {
	fmt.Fprintf(conn, "%s\n", message)
	fmt.Printf("<< %s\n", message)
}

func sendtotarget(target, message string) {
	fmt.Fprintf(conn, "PRIVMSG %s :%s\n",target, message)
	fmt.Printf("<< PRIVMSG %s :%s\n",target, message)
}

func parsefromserver(raw string) {
	fmt.Printf(">> %s", raw)

	splitted_semi := strings.Split(raw, ":")
	if(len(splitted_semi) >= 1) {
		if(len(splitted_semi[0]) > 0) {
			// REPLY TO SERVER PINGS
			if(strings.Contains(splitted_semi[0], "PING")) {
				fmt.Fprintf(conn, "PONG :%s\n", splitted_semi[1])
				fmt.Printf("<< PONG :%s\n", splitted_semi[1])
			}
		} else {
			splitted := strings.Split(splitted_semi[1], " ")

			if(len(splitted_semi) >= 1) {
				parsecommand(splitted[0] ,splitted[1], raw)
			}
		}
	}
}

func parsecommand(source, command, raw string) {
	valid := false
	for _,item := range bot.irccommands {
		if(item.command == command) {
			item.callback(source,raw)
			valid = true
		}
	}
	
	if(!valid) {
		//fmt.Printf("[!] Unknown command %s -- SOURCE %s\n", command, source)	
	}
}

/**********************************************************************************************************/

var bot = BOT{}

type BOT struct {
  irccommands []IRCCommand
  botchancommands []BotChanCommand
  botprivatecommands []BotPrivateCommand
}

func (b *BOT) addIRCcommand(command string, callback func(string,string)) {
  b.irccommands = append(b.irccommands, IRCCommand{ command: command, callback: callback })
}

func (b *BOT) addChancommand(command string, callback func(string,string,string)) {
  b.botchancommands = append(b.botchancommands, BotChanCommand{ command: command, callback: callback})
}

func (b *BOT) addPrivatecommand(command string, callback func(string,string)) {
  b.botprivatecommands = append(b.botprivatecommands, BotPrivateCommand{command: command, callback: callback})
}

type IRCCommand struct {
        command string
        callback func(string, string)
}

type BotChanCommand struct {
        command string
	callback func(string, string, string)
}

type BotPrivateCommand struct {
        command string
	callback func(string, string)
}

/**********************************************************************************************************/

func addcallbacks() {
	// IRC server stuff
	bot.addIRCcommand("PRIVMSG", callback_privmsg)

	// Public channel commands
	bot.addChancommand("!hello", callback_pub_hello)

	// Private bot commands
	bot.addPrivatecommand("help", callback_prv_help)
}

/**********************************************************************************************************/

var callback_privmsg = func(source, raw string) {
  // who is the message for ?
  // :MrSam!~sam@hobbit.r84.eu PRIVMSG #sam :jow 
  // :MrSam!~sam@hobbit.r84.eu PRIVMSG MrRobot :jowkes
  // One way or another, when splitting on " " it should be splitted[2] 
  splitted := strings.Split(raw, " ")
  target := splitted[2]

  // what is the mssage ?
  // when splitting on : it should be splitted [2]
  splitted = strings.Split(raw, ":")
  message := strings.TrimSpace(splitted[2])  

  // who is it coming from ?
  source_nick := source[0:strings.Index(source, "!")]

  // if i'm the target, get the sender's nickname 
  if(strings.EqualFold(target,my_nickname)) {
	for _,item := range bot.botprivatecommands {
		if(strings.EqualFold(item.command,message)) { item.callback(source_nick, raw) }
	}
   } else {
	for _,item := range bot.botchancommands {
                if(strings.EqualFold(item.command,message)) { item.callback(target, source_nick, raw) }
        }
  }
}

var callback_prv_help = func(source_nick, raw string) {
	sendtotarget(source_nick, "Help? I can't help " + source_nick)	
}

var callback_pub_hello = func(channel, source_nick, raw string) {
	sendtotarget(channel, "Hello " + source_nick)
}
