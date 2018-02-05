// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package sale

func init() {

	//h.MailComposeMessage().Methods().SendMail().DeclareMethod(
	//	`SendMail`,
	//	func(rs h.MailComposeMessageSet, autocommit interface{}) {
	//		//@api.multi
	//		/*def send_mail(self, auto_commit=False):
	//		  if self._context.get('default_model') == 'sale.order' and self._context.get('default_res_id') and self._context.get('mark_so_as_sent'):
	//		      order = self.env['sale.order'].browse([self._context['default_res_id']])
	//		      if order.state == 'draft':
	//		          order.state = 'sent'
	//		      self = self.with_context(mail_post_autofollow=True)
	//		  return super(MailComposeMessage, self).send_mail(auto_commit=auto_commit)
	//		*/
	//	})

}
