/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import type { ChangeEvent } from 'react'
import * as z from 'zod'
import type { Resolver } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useTranslation } from 'react-i18next'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { NativeSelect, NativeSelectOption } from '@/components/ui/native-select'
import { Switch } from '@/components/ui/switch'
import { FormDirtyIndicator } from '../components/form-dirty-indicator'
import { FormNavigationGuard } from '../components/form-navigation-guard'
import {
  SettingsForm,
  SettingsSwitchContent,
  SettingsSwitchItem,
  SettingsFormGrid,
  SettingsFormGridItem,
} from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsSection } from '../components/settings-section'
import { useSettingsForm } from '../hooks/use-settings-form'
import { useUpdateOption } from '../hooks/use-update-option'

const quotaSchema = z.object({
  QuotaForNewUser: z.coerce.number().min(0),
  PreConsumedQuota: z.coerce.number().min(0),
  TopUpLink: z.string(),
  general_setting: z.object({
    docs_link: z.string(),
  }),
  quota_setting: z.object({
    enable_free_model_pre_consume: z.boolean(),
    affiliate_reward_trigger: z.enum([
      'registration',
      'first_topup',
      'every_topup',
    ]),
    inviter_registration_reward_amount: z.coerce.number().min(0),
    invitee_registration_reward_amount: z.coerce.number().min(0),
    inviter_topup_rebate_percent: z.coerce.number().min(0),
    invitee_topup_rebate_percent: z.coerce.number().min(0),
  }),
})

type QuotaFormValues = z.infer<typeof quotaSchema>

type QuotaSettingsSectionProps = {
  defaultValues: QuotaFormValues
  complianceConfirmed?: boolean
}

export function QuotaSettingsSection({
  defaultValues,
  complianceConfirmed = true,
}: QuotaSettingsSectionProps) {
  const { t } = useTranslation()
  const updateOption = useUpdateOption()
  const handleNumberChange =
    (onChange: (value: number | string) => void) =>
    (event: ChangeEvent<HTMLInputElement>) => {
      onChange(
        event.target.value === '' ? '' : event.currentTarget.valueAsNumber
      )
    }

  const { form, handleSubmit, isDirty, isSubmitting } =
    useSettingsForm<QuotaFormValues>({
      resolver: zodResolver(quotaSchema) as Resolver<
        QuotaFormValues,
        unknown,
        QuotaFormValues
      >,
      defaultValues,
      onSubmit: async (_data, changedFields) => {
        for (const [key, value] of Object.entries(changedFields)) {
          await updateOption.mutateAsync({
            key,
            value: value as string | number | boolean,
          })
        }
      },
    })
  const rewardTrigger = form.watch('quota_setting.affiliate_reward_trigger')
  const showRegistrationRewards = rewardTrigger === 'registration'

  return (
    <SettingsSection title={t('Quota Settings')}>
      <FormNavigationGuard when={isDirty} />

      {!complianceConfirmed ? (
        <Alert variant='destructive'>
          <AlertDescription>
            {t(
              'Non-zero invitation rewards require compliance confirmation in Payment Gateway settings.'
            )}
          </AlertDescription>
        </Alert>
      ) : null}

      <Form {...form}>
        <SettingsForm onSubmit={handleSubmit}>
          <SettingsPageFormActions
            onSave={handleSubmit}
            isSaving={updateOption.isPending || isSubmitting}
          />
          <FormDirtyIndicator isDirty={isDirty} />
          <SettingsFormGrid>
            <FormField
              control={form.control}
              name='QuotaForNewUser'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('New User Quota')}</FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      value={field.value ?? ''}
                      onChange={handleNumberChange(field.onChange)}
                      name={field.name}
                      onBlur={field.onBlur}
                      ref={field.ref}
                    />
                  </FormControl>
                  <FormDescription>
                    {t('Initial quota given to new users')}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='PreConsumedQuota'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Pre-Consumed Quota')}</FormLabel>
                  <FormControl>
                    <Input
                      type='number'
                      value={field.value ?? ''}
                      onChange={handleNumberChange(field.onChange)}
                      name={field.name}
                      onBlur={field.onBlur}
                      ref={field.ref}
                    />
                  </FormControl>
                  <FormDescription>
                    {t('Quota consumed before charging users')}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='quota_setting.affiliate_reward_trigger'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Invitation reward trigger')}</FormLabel>
                  <FormControl>
                    <NativeSelect
                      value={field.value}
                      onChange={(event) => field.onChange(event.target.value)}
                      name={field.name}
                      onBlur={field.onBlur}
                      className='w-full'
                    >
                      <NativeSelectOption value='registration'>
                        {t('Registration reward')}
                      </NativeSelectOption>
                      <NativeSelectOption value='first_topup'>
                        {t('First top-up rebate')}
                      </NativeSelectOption>
                      <NativeSelectOption value='every_topup'>
                        {t('Every top-up rebate')}
                      </NativeSelectOption>
                    </NativeSelect>
                  </FormControl>
                  <FormDescription>
                    {t(
                      'Choose when referral rewards are issued to inviters and invitees.'
                    )}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            {showRegistrationRewards ? (
              <>
                <FormField
                  control={form.control}
                  name='quota_setting.inviter_registration_reward_amount'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('Inviter registration reward (USD)')}
                      </FormLabel>
                      <FormControl>
                        <Input
                          type='number'
                          min={0}
                          step='0.01'
                          value={field.value ?? ''}
                          onChange={handleNumberChange(field.onChange)}
                          name={field.name}
                          onBlur={field.onBlur}
                          ref={field.ref}
                        />
                      </FormControl>
                      <FormDescription>
                        {t(
                          'Fixed USD reward added to the inviter after a referred user registers.'
                        )}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name='quota_setting.invitee_registration_reward_amount'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>
                        {t('Invitee registration reward (USD)')}
                      </FormLabel>
                      <FormControl>
                        <Input
                          type='number'
                          min={0}
                          step='0.01'
                          value={field.value ?? ''}
                          onChange={handleNumberChange(field.onChange)}
                          name={field.name}
                          onBlur={field.onBlur}
                          ref={field.ref}
                        />
                      </FormControl>
                      <FormDescription>
                        {t(
                          'Fixed USD reward added to the invited user after registration.'
                        )}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </>
            ) : (
              <>
                <FormField
                  control={form.control}
                  name='quota_setting.inviter_topup_rebate_percent'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Inviter top-up rebate (%)')}</FormLabel>
                      <FormControl>
                        <Input
                          type='number'
                          min={0}
                          step='0.01'
                          value={field.value ?? ''}
                          onChange={handleNumberChange(field.onChange)}
                          name={field.name}
                          onBlur={field.onBlur}
                          ref={field.ref}
                        />
                      </FormControl>
                      <FormDescription>
                        {t(
                          "Percentage of the invited user's actual wallet top-up payment awarded to the inviter."
                        )}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name='quota_setting.invitee_topup_rebate_percent'
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{t('Invitee top-up rebate (%)')}</FormLabel>
                      <FormControl>
                        <Input
                          type='number'
                          min={0}
                          step='0.01'
                          value={field.value ?? ''}
                          onChange={handleNumberChange(field.onChange)}
                          name={field.name}
                          onBlur={field.onBlur}
                          ref={field.ref}
                        />
                      </FormControl>
                      <FormDescription>
                        {t(
                          "Percentage of the invited user's actual wallet top-up payment awarded back to the invitee."
                        )}
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </>
            )}

            <SettingsFormGridItem span='full'>
              <FormField
                control={form.control}
                name='quota_setting.enable_free_model_pre_consume'
                render={({ field }) => (
                  <SettingsSwitchItem>
                    <SettingsSwitchContent>
                      <FormLabel>{t('Pre-Consume for Free Models')}</FormLabel>
                      <FormDescription>
                        {t(
                          'When enabled, zero-cost models also pre-consume quota before final settlement.'
                        )}
                      </FormDescription>
                    </SettingsSwitchContent>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                        disabled={updateOption.isPending}
                      />
                    </FormControl>
                  </SettingsSwitchItem>
                )}
              />
            </SettingsFormGridItem>

            <FormField
              control={form.control}
              name='TopUpLink'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Top-Up Link')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('https://example.com/topup')}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    {t('External link for users to purchase quota')}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name='general_setting.docs_link'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('Documentation Link')}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t('https://docs.example.com')}
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    {t('Link to your documentation site')}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
          </SettingsFormGrid>
        </SettingsForm>
      </Form>
    </SettingsSection>
  )
}
